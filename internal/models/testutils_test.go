package models

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
	"testing"
)

func newTestDB(t *testing.T) *sqlx.DB {
	db, err := sqlx.Open("postgres", "user=postgres password=postgres dbname=snippetbox_test sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}

		db.Close()
	})

	return db
}
