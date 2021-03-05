package mysql_test

import (
	"os"
	"testing"
)

const envDatabaseURL = "DATABASE_URL"

var databaseURL = "api:jwjii3ke8dnsjGdsx0@tcp(10.2.1.20:3306)/billing"

func TestMain(m *testing.M) {
	if env := os.Getenv(envDatabaseURL); env != "" {
		databaseURL = env
	}
	os.Exit(m.Run())
}
