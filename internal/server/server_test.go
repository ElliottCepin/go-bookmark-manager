package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ElliottCepin/go-bookmark-manager/internal/domain"
	"github.com/ElliottCepin/go-bookmark-manager/internal/store"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"slices"
	"strconv"
	"testing"
)

var dbcount int = 0

func databaseName(tempdir string) string {
	dbcount++
	return fmt.Sprintf("%vtest-%v.db", tempdir, dbcount)
}

func equals(bm1 *domain.Bookmark, bm2 *domain.Bookmark) bool {
	return bm1.Id == bm2.Id && bm1.URL == bm2.URL && slices.Equal(bm1.Tags, bm2.Tags)
}

func TestRoundTrip(t *testing.T) {
	sqliteStore, err := store.NewSQLiteStore(databaseName(t.TempDir()))

	if err != nil {
		t.Errorf("Error setting up store: %v", err)
	}

	srv := NewServer(sqliteStore, slog.Default())
	server := httptest.NewServer(srv.Routes())

	bm := domain.Bookmark{
		URL:   "https://elliottcepin.dev",
		Title: "Portfolio Site",
		Tags: []string{
			"epic",
			"rat",
			"cow",
		},
	}

	body, err := json.Marshal(bm)

	if err != nil {
		t.Errorf("Error while marshalling: %v", err)
	}

	buf := bytes.NewBuffer(body)

	res, err := http.Post(server.URL+"/bookmarks", "application/json", buf)

	if err != nil {
		t.Errorf("Error during request: %v", err)
	}

	if res.StatusCode == http.StatusBadRequest {
		t.Errorf("Bad request on POST /bookmarks")
	}

	body, err = io.ReadAll(res.Body)

	if err != nil {
		t.Errorf("Error reading body: %v", err)
	}

	id, err := strconv.ParseInt(string(body), 10, 64)

	if err != nil {
		t.Errorf("Error parsing int: %v. Had %v, got %v", err, string(body), id)
	}

	bm.Id = id

	res, err = http.Get(server.URL + "/bookmarks/" + strconv.FormatInt(id, 10))

	if err != nil {
		t.Errorf("Error during request: %v", err)
	}

	if res.StatusCode == http.StatusBadRequest {
		t.Errorf("Bad request on GET /bookmarks/%v", id)
	}

	dec := json.NewDecoder(res.Body)

	var decodedBm domain.Bookmark
	err = dec.Decode(&decodedBm)

	if err != nil {
		t.Errorf("Error decoding JSON body: %v", err)
	}

	if !equals(&bm, &decodedBm) {
		t.Errorf("Bookmark %v:%v with %v not equal to Decoded Bookmark %v:%v with %v",
			bm.Id, bm.URL, bm.Tags, decodedBm.Id, decodedBm.URL, decodedBm.Tags)
	}

}

func TestFilter(t *testing.T) {
	sqliteStore, err := store.NewSQLiteStore(databaseName(t.TempDir()))

	if err != nil {
		t.Errorf("Error setting up store: %v", err)
	}

	srv := NewServer(sqliteStore, slog.Default())
	server := httptest.NewServer(srv.Routes())

	bm := domain.Bookmark{
		URL:   "https://elliottcepin.dev",
		Title: "Portfolio Site",
		Tags: []string{
			"epic",
			"rat",
			"cow",
		},
	}

	body, err := json.Marshal(bm)

	if err != nil {
		t.Errorf("Error while marshalling: %v", err)
	}

	buf := bytes.NewBuffer(body)

	res, err := http.Post(server.URL+"/bookmarks", "application/json", buf)

	if err != nil {
		t.Errorf("Error during request: %v", err)
	}

	if res.StatusCode == http.StatusBadRequest {
		t.Errorf("Bad request on POST /bookmarks")
	}

	body, err = io.ReadAll(res.Body)

	if err != nil {
		t.Errorf("Error reading body: %v", err)
	}

	id, err := strconv.ParseInt(string(body), 10, 64)

	if err != nil {
		t.Errorf("Error parsing int: %v. Had %v, got %v", err, string(body), id)
	}

	
	bm1 := domain.Bookmark{
		URL:   "https://coolmathgames.com/",
		Title: "Cool Math Games",
		Tags: []string{
			"epic",
			"dog",
			"cat",
		},
	}

	body, err = json.Marshal(bm1)

	if err != nil {
		t.Errorf("Error while marshalling: %v", err)
	}

	buf = bytes.NewBuffer(body)

	res, err = http.Post(server.URL+"/bookmarks", "application/json", buf)

	if err != nil {
		t.Errorf("Error during request: %v", err)
	}

	if res.StatusCode == http.StatusBadRequest {
		t.Errorf("Bad request on POST /bookmarks")
	}

	body, err = io.ReadAll(res.Body)

	if err != nil {
		t.Errorf("Error reading body: %v", err)
	}

	id1, err := strconv.ParseInt(string(body), 10, 64)

	if err != nil {
		t.Errorf("Error parsing int: %v. Had %v, got %v", err, string(body), id)
	}

	bm1.Id = id1

	res, err = http.Get(server.URL+"/bookmarks?tag=cat")

	if (err != nil) {
		t.Errorf("Error filtering by tags: %v", err)
	}

	
	if res.StatusCode == http.StatusBadRequest {
		t.Errorf("Bad request on POST /bookmarks")
	}

	var decodedBm1 []domain.Bookmark

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&decodedBm1)

	if err != nil {
		t.Errorf("Error decoding JSON: %v", err)
	}
	
	if (len(decodedBm1) != 1) {
		t.Errorf("Expected length of 1, got %v; see %v", len(decodedBm1), decodedBm1)	
	}

	if (!equals(&bm1, &decodedBm1[0])) {
		t.Errorf("Expected %v, got %v", bm, decodedBm1[0])	
	}


	
}
