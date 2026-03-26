package product

import (
	"github.com/moges7624/merkato_std/internal/filter"
)

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) GetProducts(filters *ProductFilters) ([]*Product,
	filter.Metadata, error,
) {
	products, total, err := s.store.getProducts(filters)
	if err != nil {
		return nil, filter.Metadata{}, err
	}

	metadata := filter.CalculateMetadata(total, filters.Page, filters.PageSize)
	return products, metadata, nil
}

func (s *Service) GetProduct(id int64) (*Product, error) {
	prod, err := s.store.getProduct(id)
	if err != nil {
		return nil, err
	}

	return prod, nil
}

func (s *Service) CreateProducts(cpr *CreateProductRequest) (*Product, error) {
	product := &Product{
		Name:         *cpr.Name,
		PriceInCents: int32(*cpr.PriceInDollar * 100),
		Quantity:     *cpr.Qauntity,
	}
	err := s.store.createProduct(product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *Service) DeductStock(productID int64, quantity int) error {
	p, err := s.store.getProduct(productID)
	if err != nil {
		return err
	}

	if p.Quantity < int32(quantity) {
		return ErrInsufficientStock
	}

	p.Quantity -= int32(quantity)
	return s.store.updateProduct(p)
}
