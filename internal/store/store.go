package store

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"github.com/ElliottCepin/go-bookmark-manager/internal/domain"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(filename string) (*SQLiteStore, error) {
	var err error
	s := &SQLiteStore{}

	s.db, err = sql.Open("sqlite", filename)
	
	if err != nil {
		return nil, err
	}
	
	if err := s.db.Ping(); err != nil {
		return nil, err
	}

	return s, err	
}

func (s *SQLiteStore) CreateTag(name string) (int, error) {
	return nil	
}

func (s *SQLiteStore) createBookmarkTag(tagId int, bookmarkId int) error {
	return nil
}

func (s *SQLiteStore) CreateBookmark(b *Bookmark) error {
	return nil
}

func (s *SQLiteStore) DeleteBookmark(bookmarkId int) error {
	return nil
}

func (s *SQLiteStore) FilterByBookmarkTag(tagId int) error {
	return nil
}

func (s *SQLiteStore) GetBookmark(bookmarkId int) (string, []string, error) {
	return nil
}
