package gutenberg

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	// BookSizeLimit defines maksimum size of book data that is allowed to be fetched.
	BookSizeLimit = 1024 * 1024 * 100
	// FetchTimeout defines maximum time for book fetches.
	FetchTimeout = 30 * time.Second

	gutenbergAddr = "https://www.gutenberg.org"
)

// Source can fetch books from gutenberg.org by id.
type Source struct {
	httpCli *http.Client
}

// NewSource creates new Source.
func NewSource() *Source {
	cli := http.Client{
		Timeout: FetchTimeout,
	}
	return &Source{
		httpCli: &cli,
	}
}

// Fetch fetches book with given id from gutenberg.org.
// If book is not found, returns (false, "", nil).
func (s *Source) Fetch(id string) (bool, string, error) {
	fileName, err := s.findFileName(id)
	if err != nil {
		return false, "", fmt.Errorf("searching for book txt file: %w", err)
	}
	if fileName == "" {
		return false, "", nil
	}

	found, data, err := s.fetchBookFile(id, fileName)
	if err != nil {
		return false, "", fmt.Errorf("fetching book file: %w", err)
	}
	if !found {
		return false, "", nil
	}

	return true, string(data), nil
}

func (s *Source) findFileName(id string) (string, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/files/%s/?F=0", gutenbergAddr, id),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("creating http request: %w", err)
	}

	found, listing, err := s.fetch(req)
	if err != nil {
		return "", err
	}
	if !found {
		return "", nil
	}

	fileName := regexp.MustCompile(`[a-zA-Z\d\-_]+\.txt`).FindString(listing)
	return fileName, nil
}

func (s *Source) fetchBookFile(id string, fileName string) (bool, string, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/files/%s/%s", gutenbergAddr, id, fileName),
		nil,
	)
	if err != nil {
		return false, "", fmt.Errorf("creating http request: %w", err)
	}

	return s.fetch(req)
}

func (s *Source) fetch(req *http.Request) (bool, string, error) {
	start := time.Now()
	resp, err := s.httpCli.Do(req)
	duration := time.Since(start)
	if err != nil {
		return false, "", fmt.Errorf("making http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, "", nil
	}
	if resp.StatusCode != http.StatusOK {
		return false, "", fmt.Errorf("got invalid http response status %d", resp.StatusCode)
	}

	data, err := readBody(resp.Body)
	if err != nil {
		return true, "", fmt.Errorf("reading response: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"size":     len(data),
		"respTime": duration.String(),
	}).Debugf("gutenberg: fetched url: %s", req.URL.String())

	return true, string(data), nil
}

func readBody(r io.Reader) (string, error) {
	data, err := ioutil.ReadAll(io.LimitReader(r, BookSizeLimit))
	if err != nil {
		return "", err
	}

	return string(data), nil
}
