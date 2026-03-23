package order

import (
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/moges7624/merkato_std/internal/product"
	"github.com/moges7624/merkato_std/internal/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllOrders(t *testing.T) {
	prodRepo := new(product.MockRepository)
	prodSvc := product.NewService(prodRepo)

	userRepo := new(user.MockRepository)
	userSvc := user.NewService(userRepo)

	orderRepo := new(MockRepo)
	orderSvc := NewService(orderRepo, *prodSvc, *userSvc)

	products := map[int]*product.Product{
		1: {
			ID:           1,
			Name:         gofakeit.Product().Name,
			PriceInCents: 4400,
			Quantity:     22,
		},
		2: {
			ID:           2,
			Name:         gofakeit.Product().Name,
			PriceInCents: 5500,
			Quantity:     22,
		},
	}

	userRepo.On("getUser", 1).Return(&user.User{ID: 1, Name: "john"}, nil)

	var custID int64 = 1
	var prod1Id int64 = 1
	var prod2Id int64 = 2
	var price1InUSD float32 = 33
	var price2InUSD float32 = 33

	t.Run("successful flow", func(t *testing.T) {
		item1Qt := 11
		item2Qt := 22

		params := &CreateOrderRequest{
			CustomerID: &custID,
			Items: &[]orderItemRequest{
				{
					ProductID:  &prod1Id,
					Quantity:   &item1Qt,
					PriceInUSD: &price1InUSD,
				},
				{
					ProductID:  &prod2Id,
					Quantity:   &item2Qt,
					PriceInUSD: &price2InUSD,
				},
			},
		}

		for i := range *params.Items {
			prodRepo.On("getProduct", products[i+1].ID).Return(products[i+1], nil)

			prodRepo.On("updateProduct", mock.AnythingOfType("*product.Product")).Return(nil)
		}

		orderRepo.On("insert", mock.AnythingOfType("*order.Order")).Return(nil)

		o, err := orderSvc.CreateOrder(params)

		assert.NoError(t, err)
		assert.Equal(t, o.TotalAmountInCents, int32(108900))
		assert.Len(t, o.Items, 2)

		userRepo.AssertExpectations(t)
		prodRepo.AssertExpectations(t)
		orderRepo.AssertExpectations(t)
	})

	t.Run("out of stock", func(t *testing.T) {
		item1Qt := 1
		item2Qt := 2

		params := &CreateOrderRequest{
			CustomerID: &custID,
			Items: &[]orderItemRequest{
				{
					ProductID:  &prod1Id,
					Quantity:   &item1Qt,
					PriceInUSD: &price1InUSD,
				},
				{
					ProductID:  &prod2Id,
					Quantity:   &item2Qt,
					PriceInUSD: &price2InUSD,
				},
			},
		}

		for i := range *params.Items {
			prodRepo.On("getProduct", products[i+1].ID).Return(products[1], nil)

			prodRepo.On("updateProduct", mock.AnythingOfType("*product.Product")).Return(nil)
		}

		orderRepo.On("insert", mock.AnythingOfType("*order.Order")).Return(nil)

		_, err := orderSvc.CreateOrder(params)

		assert.Error(t, err)

		userRepo.AssertExpectations(t)
		prodRepo.AssertExpectations(t)
		orderRepo.AssertExpectations(t)
	})
}
