package product

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestGetAllProducts(t *testing.T) {
	mockRepo := new(MockRepository)
	prodSvc := NewService(mockRepo)

	name1 := gofakeit.Product().Name
	name2 := gofakeit.Product().Name
	dummyProds := []*Product{
		{ID: 1, Name: name1, PriceInCents: gofakeit.Int32(), Quantity: gofakeit.Int32(), CreatedAt: gofakeit.Date()},
		{ID: 2, Name: name2, PriceInCents: gofakeit.Int32(), Quantity: gofakeit.Int32(), CreatedAt: gofakeit.Date()},
	}

	mockRepo.On("getProducts").Return(dummyProds, nil)

	prods, err := prodSvc.GetProducts()

	assert.NoError(t, err)
	assert.Len(t, prods, 2)
	assert.Equal(t, name1, prods[0].Name)

	mockRepo.AssertExpectations(t)
}

func TestGetAllProductsEmpty(t *testing.T) {
	mockRepo := new(MockRepository)
	prodSvc := NewService(mockRepo)

	mockRepo.On("getProducts").Return([]*Product{}, nil)

	prods, err := prodSvc.GetProducts()

	assert.NoError(t, err)
	assert.Empty(t, prods)
}
