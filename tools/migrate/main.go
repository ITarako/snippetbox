package main

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	databaseUrl := "postgres://postgres:postgres@localhost:5432/snippetbox?sslmode=disable"
	m, err := migrate.New("file://tools/migrate/migrations", databaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	version, dirty, err := m.Version()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Applied migration: %d, Dirty %t\n", version, dirty)

	//if err := m.Down(); err != nil && err != migrate.ErrNoChange {
	//	log.Fatal(err)
	//}
}
