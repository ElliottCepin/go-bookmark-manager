package store

import (
	"testing"
	"fmt"
	"path/filepath"
	"slices"
	"github.com/ElliottCepin/go-bookmark-manager/internal/domain"
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

func TestGetDeleted(t *testing.T) {
		
	s, err := NewSQLiteStore(filepath.Join(t.TempDir(), fmt.Sprintf("db-%v", readCount())))
		
	if (err != nil) {
		t.Errorf("Error during database setup: %v", err)
	}

	bm, err := s.CreateBookmark("https://test.site/", "Test Bookmark", []string{"Test", "Bookmark"})

	if (err != nil) {
		t.Errorf("Error creating bookmark: %v", err)
	}

	s.DeleteBookmark(bm.Id)

	bm2, err := s.GetBookmark(bm.Id)

	if (err == nil) {
		t.Errorf("Expected an error, got <nil>. returned: %v", bm2)
	}
}

func TestFilterTag(t *testing.T) {
	
	s, err := NewSQLiteStore(filepath.Join(t.TempDir(), fmt.Sprintf("db-%v", readCount())))
		
	if (err != nil) {
		t.Errorf("Error during database setup: %v", err)
	}

	bm0, err := s.CreateBookmark("https://test.site0/", "Test Bookmark 0", []string{"Test", "Bookmark"})
	bm1, err := s.CreateBookmark("https://test.site1/", "Test Bookmark 1", []string{"Test", "Bookmark"})
	bm2, err := s.CreateBookmark("https://test.site2/", "Test Bookmark 2", []string{"Rat", "Cow"})
	bm3, err := s.CreateBookmark("https://test.site3/", "Test Bookmark 3", []string{"Test", "Cow"})
	
	
	
	hasTestTag := [][]*domain.Bookmark{
		{bm0, bm1, bm3},
		{bm2},
		{bm2, bm3},
		{bm0, bm1},
		{},
	}
	tags := []string{"Test", "Rat", "Cow", "Bookmark", "Antelope"}

	for i, _ := range tags {
		tests, err := s.FilterByTag(tags[i])

		if (err != nil) {
			t.Errorf("Error during filtering '%v': %v", tags[i], err) 
		}

		if (!compareBookmarkSlice(tests, hasTestTag[i])) {
			t.Errorf("Expected '%v' filter to return %v, got %v instead", tags[i], hasTestTag, tests)
		}
	}


	
}

func compareBookmarkSlice(a []*domain.Bookmark, b []*domain.Bookmark) bool {
	if (len(a) != len(b)) {
		return false
	}
	
	matched := make([]int, 0)
	
	for _, bm1 := range a {
		for j, bm2 := range b {
			if (compareBookmarks(bm1, bm2) && !slices.Contains(matched, j)) {
				matched = append(matched, j)
				break
			}
		}
	}

	if (len(matched) == len(b)) {
		return true
	}

	return false
}

func compareBookmarks(bm1 *domain.Bookmark, bm2 *domain.Bookmark) bool {
	return (bm1.Id == bm2.Id) && (bm1.Title == bm2.Title) && (bm1.URL == bm2.URL) && (slices.Equal(bm1.Tags, bm2.Tags))	
}
