package utils

import (
	"database/sql"
	"errors"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func migrationsPath() string {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)

	return "file://" + filepath.Join(basePath, "../../migrations")
}

func NewTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres",
		"postgres://merkato:123456@localhost/test_merkatostd?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		t.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../migrations",
		"postgres",
		driver)
	if err != nil {
		t.Fatal(err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		defer db.Close()

		// run down migratoins here
		err := m.Down()
		if err != nil {
			t.Fatal(err)
		}
	})
	return db
}
