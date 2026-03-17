package product

import (
	"database/sql"
	"errors"
)

type PostgresStore struct {
	DB *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{
		DB: db,
	}
}

func (ps *PostgresStore) getProducts() (*[]Product, error) {
	query := `
	SELECT id, name, price_in_cents, quantity, created_at
	FROM products
`

	rows, err := ps.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := []Product{}

	for rows.Next() {
		var p Product

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.PriceInCents,
			&p.Quantity,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return &products, nil
}

func (ps *PostgresStore) getProduct(id int64) (*Product, error) {
	query := `
	SELECT id, name, price_in_cents, quantity, created_at
	FROM products
	WHERE id = $1
`

	var p Product
	err := ps.DB.QueryRow(query, id).Scan(
		&p.ID,
		&p.Name,
		&p.PriceInCents,
		&p.Quantity,
		&p.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	return &p, nil
}

func (ps *PostgresStore) createProduct(p *Product) error {
	query := `
	INSERT INTO products (name, price_in_cents, quantity)
	VALUES ($1, $2, $3)
	RETURNING id, created_at
	`

	err := ps.DB.QueryRow(query, p.Name, p.PriceInCents, p.Quantity).Scan(
		&p.ID,
		&p.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PostgresStore) updateProduct(p *Product) error {
	query := `
	UPDATE products
	SET name = $1, price_in_cents = $2, quantity = $3
	WHERE id = $4
`

	_, err := ps.DB.Exec(query, p.Name, p.PriceInCents, p.Quantity, p.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrProductNotFound
		}

		return err
	}

	return nil
}
