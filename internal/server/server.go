package server

import (
	"github.com/ElliottCepin/go-bookmark-manager/internal/domain"
	"log/slog"
	"net/http"
	"strings"
	"encoding/json"
)

type Store interface {
	CreateTag(string) (int64, error)
	createBookmarkTag(int64, int64) error
	CreateBookmark(string, string, []string) (*domain.Bookmark, error) // pass an object for immutability
	DeleteBookmark(int64) error 
	FilterByBookmarkTag(string) ([]*domain.Bookmark, error) // select * where id=id
	GetBookmark(int64) (*domain.Bookmark, error)
}

type Server struct {
	store Store
	logger *slog.Logger
}

func newServer(st Store, log *slog.Logger) *Server {
	s := &Server {
		store: st,
		logger: log,
	}

	return s
}

func (s *Server) createBookmark(w http.ResponseWriter, r *http.Request) {
	if (r.Method != "POST") {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	
	if (!strings.Contains(r.Header.Get("Content-Type"), "application/json")) {
		w.WriteHeader(http.StatusUnsupportedMediaType)
	}

	dec := json.NewDecoder(r.Body)
	var bm domain.Bookmark

	err := dec.Decode(&bm)
	
	if (err != nil) {
		w.WriteHeader(http.StatusBadRequest)
	}
	
	_, err = s.store.CreateBookmark(bm.URL, bm.Title, bm.Tags)
}	

