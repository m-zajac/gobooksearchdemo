package bookcache

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"

	"github.com/m-zajac/gobooksearchdemo/internal/app"
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

// BoltCache provides persistent k-v storage for book data using local fs file.
// Stores gzipped data.
// BoltCache implements apps BookCache,
type BoltCache struct {
	db         *bolt.DB
	bucketName string
}

// NewBoltCache creates new BoltCache.
func NewBoltCache(dbFilePath string) (*BoltCache, error) {
	db, err := bolt.Open(dbFilePath, 0666, nil)
	if err != nil {
		return nil, fmt.Errorf("creating bold db under file path '%s': %w", dbFilePath, err)
	}

	return &BoltCache{
		db:         db,
		bucketName: "cache",
	}, nil
}

// Fetch fetches book data from db.
func (c *BoltCache) Fetch(id string) (*app.BookCacheEntry, error) {
	var entry *app.BookCacheEntry
	err := c.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(c.bucketName))
		if bucket == nil {
			return nil
		}
		entryData := bucket.Get([]byte(id))
		if entryData == nil {
			return nil
		}
		e, err := decodeCacheEntry(entryData)
		if err != nil {
			return err
		}
		entry = e

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("fetching data from bolt bucket: %w", err)
	}

	if entry != nil {
		logrus.WithFields(logrus.Fields{
			"bookId": id,
		}).Debug("boltCache: fetched book")
	}

	return entry, nil
}

// Store stores book data in db.
func (c *BoltCache) Store(id string, entry app.BookCacheEntry) error {
	entryData, err := encodeCacheEntry(entry)
	if err != nil {
		return err
	}

	err = c.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(c.bucketName))
		if err != nil {
			return fmt.Errorf("creating bolt bucket: %w", err)
		}
		if err = bucket.Put([]byte(id), entryData); err != nil {
			return fmt.Errorf("putting data to bolt bucket: %w", err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"bookId":    id,
		"entrySize": len(entryData),
	}).Info("boltCache: stored book")

	return nil
}

// Close closes underlying bolt db.
func (c *BoltCache) Close() error {
	if err := c.db.Close(); err != nil {
		return fmt.Errorf("closing bolt db: %w", err)
	}
	return nil
}

func encodeCacheEntry(entry app.BookCacheEntry) ([]byte, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	if err := json.NewEncoder(zw).Encode(entry); err != nil {
		return nil, fmt.Errorf("marshaling cache entry to json: %w", err)
	}

	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("closing gzip writer: %w", err)
	}

	return buf.Bytes(), nil
}

func decodeCacheEntry(b []byte) (*app.BookCacheEntry, error) {
	buf := bytes.NewBuffer(b)
	zr, err := gzip.NewReader(buf)
	if err != nil {
		return nil, fmt.Errorf("creating gzip reader: %w", err)
	}

	var e app.BookCacheEntry
	if err := json.NewDecoder(zr).Decode(&e); err != nil {
		return nil, fmt.Errorf("decoding cache entry: %w", err)
	}

	return &e, nil
}
