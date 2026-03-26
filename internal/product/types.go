package product

import (
	"github.com/moges7624/merkato_std/internal/filter"
	"github.com/moges7624/merkato_std/internal/validator"
)

type Store interface {
	getProducts(filters *ProductFilters) ([]*Product, int, error)
	getProduct(id int64) (*Product, error)
	createProduct(prod *Product) error
	updateProduct(prod *Product) error
}

type CreateProductRequest struct {
	Name          *string  `json:"name"`
	PriceInDollar *float32 `json:"price_in_dollar"`
	Qauntity      *int32   `json:"quantity"`
}

type ProductFilters struct {
	Name string
	filter.Filters
}

func ValidateCreateProductParams(v *validator.Validator, cpr CreateProductRequest) {
	v.Check(cpr.Name != nil, "name", "must be provided")
	v.Check(cpr.PriceInDollar != nil, "price_in_dollar", "must be provided")
	if cpr.PriceInDollar != nil && *cpr.PriceInDollar < 0 {
		v.AddError("price_in_dollar", "must be greater than 0")
	}
	v.Check(cpr.Qauntity != nil, "quantity", "must be provided")
	if cpr.Qauntity != nil && *cpr.Qauntity < 0 {
		v.AddError("quantity", "must be greater than 0")
	}
}
