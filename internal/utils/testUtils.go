package utils

import (
	"database/sql"
	"errors"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/moges7624/merkato_std/internal/product"
	"golang.org/x/crypto/bcrypt"
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

		err := m.Down()
		if err != nil {
			t.Fatal(err)
		}
	})
	return db
}

func SeedDB(t *testing.T, db *sql.DB, table string) {
	switch table {
	case "users":
		seedUser(t, db)
	case "products":
		seedProduct(t, db)
	default:
		t.Fatal("SeedDB: invalid table name")
	}
}

func seedUser(t *testing.T, db *sql.DB) {
	query := `
	INSERT INTO users (name, email, password_hash)
	VALUES ($1, $2, $3)
	`

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("pass1234"), 12)
	if err != nil {
		t.Fatalf("error hashing password, %v", err)
	}

	_, err = db.Exec(query, "Adams", "adams@mail.com", passwordHash)
	if err != nil {
		t.Fatalf("error seeding user, %v", err)
	}
}

func seedProduct(t *testing.T, db *sql.DB) {
	p := product.Product{
		Name:         gofakeit.Product().Name,
		PriceInCents: 8800,
		Quantity:     34,
	}

	query := `
	INSERT INTO products (name, price_in_cents, quantity)
	VALUES ($1, $2, $3)
	RETURNING id, created_at
	`

	_, err := db.Exec(query, p.Name, p.PriceInCents, p.Quantity)
	if err != nil {
		t.Fatalf("error seeding product, %v", err)
	}
}
