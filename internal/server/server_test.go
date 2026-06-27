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

func databaseName() string {
	dbcount++
	return fmt.Sprintf("test-%v.db", dbcount)
}

func equals(bm1 *domain.Bookmark, bm2 *domain.Bookmark) bool {
	return bm1.Id == bm2.Id && bm1.URL == bm2.URL && slices.Equal(bm1.Tags, bm2.Tags)
}

func TestRoundTrip(t *testing.T) {
	sqliteStore, err := store.NewSQLiteStore(databaseName())

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
