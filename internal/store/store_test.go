package store

import (
	"testing"
	"fmt"
	"path/filepath"
)

var count int = 0

func readCount() int {
	count++
	return count
}

func TestConnectDB(t *testing.T) {
	_, err := NewSQLiteStore(filepath.Join(t.TempDir(), fmt.Sprintf("db-%v", readCount())))
	
	if err != nil {
		t.Errorf("Error during database setup: %v", err)
	}
}

func TestRoundTrip(t *testing.T) {
	s, err := NewSQLiteStore(filepath.Join(t.TempDir(), fmt.Sprintf("db-%v", readCount())))
	
	if (err != nil) {
		t.Errorf("Error during database setup: %v", err)
	}

	bm, err = s.CreateBookmark("https://test.site/", "Test Bookmark", []string{"Test", "Bookmark"})

	if (err != nil) {
		t.Errorf("Error creating bookmark: %v", err)
	}
	
	n, t, err := s.GetBookmark(bm)
}
