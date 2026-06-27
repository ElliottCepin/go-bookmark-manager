package store

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"github.com/ElliottCepin/go-bookmark-manager/internal/domain"
	"slices"
	"errors"
	"time"
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
		name TEXT UNIQUE NOT NULL
	);`)

	if (err != nil) {
		return nil, err
	}

	_, err = s.db.Exec(`CREATE TABLE IF NOT EXISTS bookmark_tags (
		bm_id INTEGER NOT NULL,
		tag_id INTEGER NOT NULL,
		PRIMARY KEY (bm_id, tag_id)
	);`)

	if (err != nil) {
		return nil, err
	}

	return s, err	
}

func (s *SQLiteStore) CreateTag(name string) (int64, error) {
	r, err := s.db.Exec("INSERT INTO tags (name) VALUES (?);", name)
	
	if (err != nil) {
		return -1, err
	}

	id, err := r.LastInsertId()	

	if (err != nil) {
		return -1, err
	}

	return id, nil
}

func (s *SQLiteStore) CreateBookmarkTag(bmId int64, tagId int64) error {
	_, err := s.db.Exec("INSERT INTO bookmark_tags (bm_id, tag_id) VALUES (?, ?)", bmId, tagId)
	return err
}

func (s *SQLiteStore) CreateBookmark(url string, title string, tags []string) (*domain.Bookmark, error) {
	// check that tag exists then create tag
	tagIds := make([]int64, 0)
	tagNames := make([]string, 0)

	var id int64
	var name string

	for _, tag := range tags {
		row := s.db.QueryRow("SELECT id, name FROM tags WHERE name=(?)", tag)


		err := row.Scan(&id, &name)

		if (errors.Is(err, sql.ErrNoRows)) {
			name = tag
			id, err = s.CreateTag(name)

			if (err != nil) {
				return nil, err
			}
		} else if (err != nil) {
			return nil, err
		}
		
		
		if (slices.Contains(tagNames, tag)) {
			return nil, errors.New("Tags violates unique constraint: too many tags")
		}
		
		tagIds = append(tagIds, id)
		tagNames = append(tagNames, name)
	}

	r, err := s.db.Exec("INSERT INTO bookmarks (url, title) VALUES (?, ?)", url, title)

	if (err != nil) {
		return nil, err
	}

	bmid, err := r.LastInsertId()
	
	if (err != nil) {
		return nil, err
	}
	
	for _, tagId := range tagIds {
		err := s.CreateBookmarkTag(bmid, tagId)

		if (err != nil) {
			return nil, err
		}
	}

	bm := &domain.Bookmark{
		Id: bmid,
		URL: url,
		Title: title,
		Tags: tagNames,
		Time: time.Now(),
	}

	// finally, rewrite with hot queries
	// do all of these need to have a context?
	return bm, nil
}

func (s *SQLiteStore) DeleteBookmark(bookmarkId int64) error {
	_, err := s.db.Exec("DELETE FROM bookmarks where bookmarks.id=(?);", bookmarkId)
	if (err != nil) {
		return err
	}
	return nil
}

func (s *SQLiteStore) FilterByTag(tagName string) ([]*domain.Bookmark, error) {
	rows, err := s.db.Query("SELECT bm_id FROM bookmark_tags JOIN tags ON tags.name=(?) AND tags.id=bookmark_tags.tag_id;", tagName)

	if (err != nil) {
		return nil, err
	}
	
	var bm_id int64
	
	bm_ids := make([]int64, 0)

	for rows.Next() {
		err = rows.Scan(&bm_id)		

		if (err != nil) { 
			return nil, err
		}
		
		bm_ids = append(bm_ids, bm_id) // these are unique, so no error checking needed
	}
	

	bms := make([]*domain.Bookmark, 0)

	for _, bm := range bm_ids {
		bmObj, err := s.GetBookmark(bm) // to get the bookmark as an object

		if err != nil {
			return nil, err
		}

		bms = append(bms, bmObj)
	}

	return bms, nil
}

func (s *SQLiteStore) GetBookmark(bookmarkId int64) (*domain.Bookmark, error) {
	row := s.db.QueryRow("SELECT * FROM bookmarks WHERE bookmarks.id=(?);", bookmarkId)

	var id int64
	var url string
	var title string
	var createdAt time.Time

	err := row.Scan(&id, &url, &title, &createdAt)

	if (err != nil) {
		return nil, err
	}

	rows, err := s.db.Query("SELECT * FROM bookmark_tags WHERE bookmark_tags.bm_id=(?);", bookmarkId)
	defer rows.Close()
	

	tagNames := make([]string, 0)
	tagIds := make([]int64, 0)
	var bmId int64
	var tagId int64
	var tagName string

	for rows.Next() {

		err = rows.Scan(&bmId, &tagId)

		if err != nil {
			return nil, err
		}

		tagIds = append(tagIds, tagId)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}


	var nullInt int64

	for _, tagId := range tagIds { 
		row = s.db.QueryRow("SELECT id, name FROM tags WHERE id=(?)", tagId)

		if err = row.Scan(&nullInt, &tagName); err != nil {
			return nil, err
		}

		tagNames = append(tagNames, tagName)	
	}

	
	bm := &domain.Bookmark{
		Id: id,
		URL: url,
		Title: title,
		Tags: tagNames,
		Time: createdAt,
	}

	return bm, nil
}
