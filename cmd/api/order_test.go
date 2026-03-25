//go:build integration

package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/moges7624/merkato_std/internal/assert"
	"github.com/moges7624/merkato_std/internal/order"
	"github.com/moges7624/merkato_std/internal/product"
	"github.com/moges7624/merkato_std/internal/user"
	"github.com/moges7624/merkato_std/internal/utils"
)

func newOrderHandler(t *testing.T) *orderHandler {
	db := utils.NewTestDB(t)
	api := &APIServer{
		logger: slog.New(slog.DiscardHandler),
		DB:     db,
	}

	userRepo := user.NewPostgresStore(db)
	userSvc := user.NewService(userRepo)
	prodRepo := product.NewPostgresStore(db)
	prodSvc := product.NewService(prodRepo)
	orderRepo := order.NewPostgresStore(db)
	orderSvc := order.NewService(orderRepo, *prodSvc, *userSvc)

	return NewOrderHandler(api, *orderSvc)
}

func TestOrderHandler_CreateOrder(t *testing.T) {
	orderHandler := newOrderHandler(t)
	utils.SeedDB(t, orderHandler.s.DB, "users")
	utils.SeedDB(t, orderHandler.s.DB, "products")

	t.Run("given product with insufficient stock, it should return out of stock message",
		func(t *testing.T) {
			reqBody := `{
				"customer_id": 1,
				"items": [
				{
				"product_id": 1,
				"quantity": 100,
				"price_in_usd": 29.33
				}
				]
				}`

			req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			orderHandler.handleCreateOrder(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, http.StatusUnprocessableEntity)

			var apiRes struct {
				Error APIError
			}

			err = json.NewDecoder(res.Body).Decode(&apiRes)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, apiRes.Error.Type, InventoryError)
			assert.Equal(t, apiRes.Error.Message, product.ErrInsufficientStock.Error())
		})

	t.Run("given valid input, it should create an order", func(t *testing.T) {
		reqBody := `{
				"customer_id": 1,
				"items": [
				{
				"product_id": 1,
				"quantity": 10,
				"price_in_usd": 29.33
				}
				]
				}`

		req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		orderHandler.handleCreateOrder(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, res.StatusCode, http.StatusCreated)

		var apiRes struct {
			Order order.Order
		}

		err = json.NewDecoder(res.Body).Decode(&apiRes)
		if err != nil {
			t.Fatalf("error decoding response, %v", err)
		}

		assert.Equal(t, apiRes.Order.TotalAmountInCents, 29330)
		assert.Equal(t, apiRes.Order.Status, order.StatusPending)
	})
}
