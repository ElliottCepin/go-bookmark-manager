package store

import (
	"testing"
	"fmt"
	"path/filepath"
	"slices"
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

	bm, err := s.CreateBookmark("https://test.site/", "Test Bookmark", []string{"Test", "Bookmark"})

	if (err != nil) {
		t.Errorf("Error creating bookmark: %v", err)
	}
	
	bm2, err := s.GetBookmark(bm.Id)

	if (err != nil) {
		t.Errorf("Error retrieving Bookmark: %v", err)
	}

	if (bm2.Id != bm.Id) {
		t.Errorf("Original id '%v' does not match returned id '%v'", bm.Id, bm2.Id)	
	}

	if (bm.URL != bm2.URL) {
		t.Errorf("Original url '%v' des not match returned url '%v'", bm.URL, bm2.URL)
	}

	if (bm.Title != bm2.Title) {
		t.Errorf("Original title '%v' does not match returned title '%v'", bm.Title, bm2.Title)	
	}

	if (!slices.Equal(bm.Tags, bm2.Tags)) {
		t.Errorf("Original tags '%v' do not match returned tags '%v'", bm.Tags, bm2.Tags)	
	}
	
}
