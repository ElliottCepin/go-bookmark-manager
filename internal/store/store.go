package store

import (
	"database/sql"
	_ "modernc.org/sqlite"
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

func (s *SQLiteStore) CreateTag(name string) error {
	return nil	
}

func (s *SQLiteStore) createBookmarkTag(tagId int, bookmarkId int) error {
	return nil
}

func (s *SQLiteStore) CreateBookmark(url string, title string, tags []string) error {
	return nil
}

func (s *SQLiteStore) DeleteBookmark(bookmarkId int) error {
	return nil
}

func (s *SQLiteStore) FilterByBookmarkTag(tagId int) error {
	return nil
}

