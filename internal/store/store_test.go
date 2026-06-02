package store

import (
	"testing"
	"fmt"
	"path/filepath"
)

var count int = 0

func readCount() int {
	count++
	return count
}

func TestConnectDB(t *testing.T) {
	store := NewSQLiteStore(filepath.Join(t.TempDir(), fmt.Sprintf("db-%v", readCount())))
	
	if err := store.db.Ping(); err != nil {
		t.Errorf("Error connecting to database: %v", err)
	}


}
