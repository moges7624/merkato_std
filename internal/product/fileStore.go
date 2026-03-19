package product

import (
	"sync"

	"github.com/brianvoe/gofakeit/v7"
)

type FileStore struct {
	mu       sync.RWMutex
	products map[int64]Product
	idSeq    int64
}

func NewFileStore() *FileStore {
	return &FileStore{
		products: map[int64]Product{
			1: {
				ID:           1,
				Name:         gofakeit.Product().Name,
				PriceInCents: int32(gofakeit.Product().Price),
				Quantity:     int32(gofakeit.Number(1, 300)),
				CreatedAt:    gofakeit.Date(),
			},
			2: {
				ID:           2,
				Name:         gofakeit.Product().Name,
				PriceInCents: int32(gofakeit.Product().Price),
				Quantity:     int32(gofakeit.Number(1, 300)),
				CreatedAt:    gofakeit.Date(),
			},
		},
		idSeq: 2,
	}
}

func (s *FileStore) getProducts() ([]*Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	products := make([]*Product, 0, len(s.products))
	for _, product := range s.products {
		products = append(products, &product)
	}

	return products, nil
}

func (s *FileStore) getProduct(id int64) (*Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	prod, ok := s.products[id]
	if !ok {
		return nil, ErrProductNotFound
	}

	return &prod, nil
}

func (s *FileStore) createProduct(prod *Product) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.idSeq++
	prod.ID = s.idSeq
	s.products[prod.ID] = *prod

	return nil
}

func (s *FileStore) updateProduct(prod *Product) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.products[prod.ID] = *prod

	return nil
}
