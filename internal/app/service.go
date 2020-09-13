package app

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// BookSource returns a book content for given book id
type BookSource interface {
	Fetch(id string) (bool, string, error)
}

// BookCacheEntry defines data for storing in cache.
type BookCacheEntry struct {
	Found bool
	Data  string
}

// BookCache stores book data.
type BookCache interface {
	Fetch(id string) (*BookCacheEntry, error)
	Store(id string, entry BookCacheEntry) error
}

// Service provides method for finding paragraphs matching a phrase in books.
type Service struct {
	source BookSource
	cache  BookCache
}

// NewService creates new Service.
func NewService(bs BookSource, bc BookCache) (*Service, error) {
	if bs == nil {
		return nil, errors.New("book source is nil")
	}
	if bc == nil {
		return nil, errors.New("book cache is nil")
	}
	return &Service{
		source: bs,
		cache:  bc,
	}, nil
}

// FindBookParagraph finds a book with given id and tries to find given phrase using fuzzy search.
// Fuzziness parameter defines allowed difference level in matched text.
// If book is not found, ("", ErrBookNotFound) is returned.
// If phrase is not found in a book, ("", nil) is returned.
func (s *Service) FindBookParagraph(bookID string, phrase string, fuzziness uint) (string, error) {
	found, data, err := s.getBookData(bookID)
	if err != nil {
		return "", err
	}
	if !found {
		return "", ErrBookNotFound
	}

	lowerPhrase := strings.ToLower(phrase)
	lowerData := strings.ToLower(data)

	start := time.Now()
	phraseFound, matchIdx := FuzzySearch(
		[]rune(lowerPhrase),
		[]rune(lowerData),
		fuzziness,
	)
	duration := time.Since(start)

	logrus.Debugf("fuzzysearch: search time: %s", duration.String())

	if !phraseFound {
		return "", nil
	}

	paragraph, err := ExtractParagraph([]rune(data), matchIdx, 1000)
	if err != nil {
		return "", fmt.Errorf("extracting paragraph from book text: %w", err)
	}

	return paragraph, nil
}

func (s *Service) getBookData(bookID string) (bool, string, error) {
	cacheEntry, err := s.cache.Fetch(bookID)
	if err != nil {
		return false, "", fmt.Errorf("fetching book from cache: %w", err)
	}
	if cacheEntry != nil {
		return cacheEntry.Found, cacheEntry.Data, nil
	}

	found, data, err := s.source.Fetch(bookID)
	if err != nil {
		return true, "", fmt.Errorf("fetching book from source: %w", err)
	}
	if err := s.cache.Store(bookID, BookCacheEntry{
		Found: found,
		Data:  data,
	}); err != nil {
		return true, data, fmt.Errorf("storing book in cache: %w", err)
	}
	return found, data, nil
}
