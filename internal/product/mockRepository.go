package product

import (
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func NewMockRepo() *MockRepository {
	return &MockRepository{}
}

func (m *MockRepository) getProducts(fiters *ProductFilters) ([]*Product, int, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}

	return args.Get(0).([]*Product), args.Int(1), args.Error(2)
}

func (m *MockRepository) getProduct(id int64) (*Product, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Product), args.Error(1)
}

func (m *MockRepository) createProduct(p *Product) error {
	args := m.Called(p)

	return args.Error(0)
}

func (m *MockRepository) updateProduct(p *Product) error {
	args := m.Called(p)

	return args.Error(0)
}
