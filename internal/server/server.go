package server

import (
	"encoding/json"
	"github.com/ElliottCepin/go-bookmark-manager/internal/domain"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"fmt"
)

type Store interface {
	CreateTag(string) (int64, error)
	CreateBookmarkTag(int64, int64) error
	CreateBookmark(string, string, []string) (*domain.Bookmark, error) // pass an object for immutability
	DeleteBookmark(int64) error
	FilterByTag(string) ([]*domain.Bookmark, error) // select * where id=id
	GetBookmark(int64) (*domain.Bookmark, error)
}

type Server struct {
	store  Store
	logger *slog.Logger
}

func NewServer(st Store, log *slog.Logger) *Server {
	s := &Server{
		store:  st,
		logger: log,
	}

	return s
}

func (s *Server) createBookmark(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(r.Body)
	var bm domain.Bookmark

	err := dec.Decode(&bm)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	

	bmCopy, err := s.store.CreateBookmark(bm.URL, bm.Title, bm.Tags)

	fmt.Fprintf(w, "%v", bmCopy.Id)
}

func (s *Server) filterBookmarks(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()

	tag, ok := query["tag"]
	if !ok {
		// return empty json array
		return
	}
	
	bms, err := s.store.FilterByTag(tag[0])
	
	if (err != nil) {
		w.WriteHeader(http.StatusBadRequest)
	}

	enc := json.NewEncoder(w)

	err = enc.Encode(bms)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) deleteBookmark(w http.ResponseWriter, r *http.Request) {
	// get slug from {slug}
	sid := r.PathValue("id")
	id, err := strconv.ParseInt(sid, 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.store.DeleteBookmark(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) retrieveBookmark(w http.ResponseWriter, r *http.Request) {
	// get slug from {slug}
	sid := r.PathValue("id")
	id, err := strconv.ParseInt(sid, 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bm, err := s.store.GetBookmark(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(bm)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (s *Server) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /bookmarks", s.createBookmark)
	mux.HandleFunc("GET /bookmarks", s.filterBookmarks)
	mux.HandleFunc("GET /bookmarks/{id}", s.retrieveBookmark)
	mux.HandleFunc("DELETE /bookmarks/{id}", s.deleteBookmark)

	return mux
}
