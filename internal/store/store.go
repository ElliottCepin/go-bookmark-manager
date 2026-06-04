package store

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"github.com/ElliottCepin/go-bookmark-manager/internal/domain"
)

type SQLiteStore struct {
	db *sql.DB
	createBM *sql.Stmt
	createTag *sql.Stmt
	createBMTag *sql.Stmt
	deleteBM *sql.Stmt
	getBM *sql.Stmt
	filterByTag *sql.Stmt
}

func NewSQLiteStore(filename string) (*SQLiteStore, error) {
	var err error
	s := &SQLiteStore{}

	s.db, err = sql.Open("sqlite", filename)
	
	if err != nil {
		return nil, err
	}
	
	// this seems silly since concurrency is safe
	s.db.SetMaxOpenConns(1)
	
	if err := s.db.Ping(); err != nil {
		return nil, err
	} 

	_, err = s.db.Exec(`CREATE TABLE IF NOT EXISTS bookmarks (
		id INTEGER PRIMARY KEY,
		url TEXT NOT NULL,
		title TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP

	);`)

	if (err != nil) {
		return nil, err
	}
	
	_, err = s.db.Exec(`CREATE TABLE IF NOT EXISTS tags (
		id INTEGER PRIMARY KEY,
		name UNIQUE TEXT NOT NULL
	);`)

	if (err != nil) {
		return nil, err
	}

	_, err = s.db.Exec(`CREATE TABLE IF NOT EXISTS bookmark_tags (
		bm_id INTEGER NOT NULL,
		tag_id INTEGER NOT NULL
		PRIMARY KEY (bm_id, tag_id)
	);`)

	if (err != nil) {
		return nil, err
	}

	return s, err	
}

func (s *SQLiteStore) CreateTag(name string) (int, error) {
	return -1, nil
}

func (s *SQLiteStore) createBookmarkTag(tagId int, bookmarkId int) error {
	return nil
}

func (s *SQLiteStore) CreateBookmark(url string, title string, tags []string) (*domain.Bookmark, error) {
	// check that tag exists then create tag
	tagIds := make([]int)
	tagNames := make([]string)

	var id int
	var name string

	for _, tag := range tags {
		row := s.db.QueryRow("SELECT id, name FROM tags where name=(?)", tag)


		err := row.Scan(&id, &name)

		if (err != null) {
			return nil, err
		}
		
		
		if (slices.contains(tagNames) {
			return nil, errors.New("Tags violates unique constraint: too many tags")
		}
		
		tagIds := append(tagIds, id)
		tagNames := append(tagNames, name)
	}

	r, err := s.db.Exec("INSERT INTO bookmarks (url, title) VALUES (?, ?)", url, title)

	if (err != nil) {
		return nil, err
	}

	id, err = r.LastInsertId()
	
	if (err != nil) {
		return nil, err
	}
	
	for _, tagId := range tagIds {
		_, err := s.db.Exec("INSERT INTO bookmark_tags (bm_id, tag_id) VALUES (?, ?)", id, tagId)

		if (err != nil) {
			return nil, err
		}
	}

	bm := &domain.Bookmark{
		Id: id,
		URL: url,
		Title: title,
		Tags: tagNames
		Time: time.Now()
	}

	// finally, rewrite with hot queries
	// do all of these need to have a context?
	return bm, nil
}

func (s *SQLiteStore) DeleteBookmark(bookmarkId int) error {
	return nil
}

func (s *SQLiteStore) FilterByBookmarkTag(tagId int) error {
	return nil
}

func (s *SQLiteStore) GetBookmark(bookmarkId int) (*domain.Bookmark, error) {
	return nil, nil
}
