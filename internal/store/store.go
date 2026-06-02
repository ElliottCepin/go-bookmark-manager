package store

import (
	"database/sql"
	"modernc.org/sqlite"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(filename string) *SQLiteStore {
	//s := &SQLiteStore{}
	//s.db, err := sql.Open("sqlite3", filename)
	
	
}

func (*s SQLiteStore) CreateTag(name string) error {
	
}

func (*s SQLiteStore) createBookmarkTag(tagId int, bookmarkId int) error {

}

func (*s SQLiteStore) CreateBookmark(url string, title string, tags []string) error {

}

func (*s SQLiteStore) DeleteBookmark(bookmarkId int) error {

}

func (*s SQLiteStore) FilterByBookmarkTag(tagId int) error {

}

