package order

import (
	"database/sql"
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

func (ps *PostgresStore) getAll() ([]*Order, error) {
	query := `
	SELECT
		orders.id,
		orders.customer_id,
		orders.status,
		orders.total_amount_in_cents,
		order_items.id,
		order_items.product_id,
		order_items.quantity,
		order_items.purchase_price_in_cents
	FROM orders
	LEFT JOIN order_items 
	ON order_items.order_id = orders.id
`

	rows, err := ps.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var orders []*Order
	orderMap := make(map[int64]*Order)

	for rows.Next() {
		var o Order
		var productID, itemID sql.NullInt64
		var itemQuantity sql.NullInt16
		var purchasePriceInCents sql.NullInt32

		err := rows.Scan(
			&o.ID,
			&o.CustomerID,
			&o.Status,
			&o.TotalAmountInCents,
			&itemID,
			&productID,
			&itemQuantity,
			&purchasePriceInCents)
		if err != nil {
			return nil, err
		}

		if _, exists := orderMap[o.ID]; !exists {
			orderMap[o.ID] = &o
			orders = append(orders, &o)
		}

		if itemID.Valid {
			item := OrderItem{
				ID:                   itemID.Int64,
				ProductID:            productID.Int64,
				Quantity:             int(itemQuantity.Int16),
				PurchasePriceInCents: purchasePriceInCents.Int32,
			}

			orderMap[o.ID].Items = append(orderMap[o.ID].Items, item)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (ps *PostgresStore) getByID(id int64) (*Order, error) {
	query := `
	SELECT
		orders.id,
		orders.customer_id,
		orders.status,
		orders.total_amount_in_cents,
		order_items.product_id,
		order_items.quantity,
		order_items.purchase_price_in_cents
	FROM orders
	LEFT JOIN order_items 
	ON order_items.order_id = orders.id
	WHERE orders.id = $1
`

	var order *Order

	rows, err := ps.DB.Query(query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var o Order
		var productID sql.NullInt64
		var itemQuantity sql.NullInt16
		var purchasePriceInCents sql.NullInt32

		err := rows.Scan(
			&o.ID,
			&o.CustomerID,
			&o.Status,
			&o.TotalAmountInCents,
			&productID,
			&itemQuantity,
			&purchasePriceInCents)
		if err != nil {
			return nil, err
		}

		if order == nil {
			order = &Order{
				ID:                 o.ID,
				CustomerID:         o.CustomerID,
				Status:             o.Status,
				TotalAmountInCents: o.TotalAmountInCents,
			}
		}

		if productID.Valid {
			order.Items = append(order.Items, OrderItem{
				ProductID:            productID.Int64,
				Quantity:             int(itemQuantity.Int16),
				PurchasePriceInCents: purchasePriceInCents.Int32,
			})
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if order == nil {
		return nil, ErrOrderNotFound
	}

	return order, nil
}

func (ps *PostgresStore) insert(o *Order) error {
	tx, err := ps.DB.Begin()
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}

	defer tx.Rollback()

	query := `
	INSERT INTO orders (customer_id, status, total_amount_in_cents)
	VALUES ($1, $2, $3)
	RETURNING id
`

	var orderID int64
	err = tx.QueryRow(query, o.CustomerID, o.Status, o.TotalAmountInCents).Scan(&orderID)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	o.ID = orderID

	itemQuery := `
	INSERT INTO order_items (order_id, product_id, quantity, purchase_price_in_cents)
	VALUES ($1, $2, $3, $4)
`
	for _, item := range o.Items {
		args := []any{orderID, item.ProductID, item.Quantity, item.PurchasePriceInCents}
		_, err := tx.Exec(itemQuery, args...)
		if err != nil {
			return fmt.Errorf("failed to insert order item: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}
