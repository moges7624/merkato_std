package product

import (
	"database/sql"
	"errors"
	"fmt"
)

type PostgresStore struct {
	DB *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{
		DB: db,
	}
}

func (ps *PostgresStore) getProducts(filters *ProductFilters) ([]*Product, int, error) {
	query := fmt.Sprintf(
		`SELECT count(*) OVER(), id, name, price_in_cents, quantity, created_at
	FROM products
	WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
	ORDER BY %s %s, id ASC
	LIMIT $2 OFFSET $3`,
		filters.SortColumn(),
		filters.SortDirection())

	rows, err := ps.DB.Query(query, filters.Name, filters.Limit(), filters.Offset())
	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}

	defer rows.Close()

	totalRecords := 0
	products := []*Product{}

	for rows.Next() {
		var p Product

		err := rows.Scan(
			&totalRecords,
			&p.ID,
			&p.Name,
			&p.PriceInCents,
			&p.Quantity,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		products = append(products, &p)
	}

	return products, totalRecords, nil
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
