package mysql

import (
	"database/sql"
	"testing"
)

func TestStore(t *testing.T, databaseURL string) *Store {
	t.Helper()

	db, err := sql.Open("mysql", databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	return New(db)
}
