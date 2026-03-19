package order

import (
	"fmt"

	"github.com/moges7624/merkato_std/internal/validator"
)

type Store interface {
	getAll() ([]*Order, error)
	insert(order *Order) error
	getByID(id int64) (*Order, error)
}

type orderItemRequest struct {
	ProductID  *int64   `json:"product_id"`
	Quantity   *int     `json:"quantity"`
	PriceInUSD *float32 `json:"price_in_usd"`
}

type CreateOrderRequest struct {
	CustomerID *int64              `json:"customer_id"`
	Items      *[]orderItemRequest `json:"items"`
}

func ValidateOrderRequestItmes(
	v *validator.Validator,
	items *[]orderItemRequest,
) {
	for i, item := range *items {
		v.Check(
			item.ProductID != nil,
			fmt.Sprintf("items-%d-product_id", i),
			"must be provided",
		)

		v.Check(
			item.Quantity != nil,
			fmt.Sprintf("items-%d-quantity", i),
			"must be provided",
		)

		if item.Quantity != nil {
			v.Check(
				*item.Quantity > 0,
				fmt.Sprintf("items-%d-quantity", i),
				"must be greater than 0",
			)
		}

		v.Check(
			item.PriceInUSD != nil,
			fmt.Sprintf("items-%d-price_in_usd", i),
			"must be provided",
		)

		if item.PriceInUSD != nil {
			v.Check(
				*item.PriceInUSD > 0,
				fmt.Sprintf("items-%d-price_in_usd", i),
				"must be greater than 0",
			)
		}
	}
}

func ValidateCreateOrderRequest(
	v *validator.Validator,
	body *CreateOrderRequest,
) {
	v.Check(body.CustomerID != nil, "user_id", "must be provided")
	v.Check(body.Items != nil, "items", "must be provided")
	v.Check(len(*body.Items) > 0, "items", "must have atleast one item")
	ValidateOrderRequestItmes(v, body.Items)
}
