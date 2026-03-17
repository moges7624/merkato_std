package main

import (
	"io"
	"net/http"

	"github.com/moges7624/merkato_std/internal/order"
	"github.com/moges7624/merkato_std/internal/product"
	"github.com/moges7624/merkato_std/internal/user"
)

func homeHandler(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "Hello there")
}

func (s *APIServer) NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/home", homeHandler)

	// userFileStore := user.NewFileStore()
	userPgStore := user.NewPostgresStore(s.DB)
	userService := user.NewService(userPgStore)
	userHandler := NewUserHandler(s, *userService)
	mux.HandleFunc("GET /users", userHandler.handleGetUsers)
	mux.HandleFunc("POST /users", userHandler.handleCreateUser)
	mux.HandleFunc("GET /users/{id}", userHandler.handleGetUser)
	mux.HandleFunc("PATCH /users/{id}", userHandler.handleUpdateUser)
	mux.HandleFunc("DELETE /users/{id}", userHandler.handleDeleteUser)

	// productFileStore := product.NewFileStore()
	productPostgresStore := product.NewPostgresStore(s.DB)
	productService := product.NewService(productPostgresStore)
	productHandler := NewProductHandler(s, *productService)
	mux.HandleFunc("GET /products", productHandler.handleGetProducts)
	mux.HandleFunc("GET /products/{id}", productHandler.handleGetProduct)
	mux.HandleFunc("POST /products", productHandler.handleCreateProduct)

	orderFileStore := order.NewFileStore()
	orderService := order.NewService(orderFileStore, *productService)
	orderHandler := NewOrderHandler(s, *orderService)
	mux.HandleFunc("GET /orders", orderHandler.handleGetOrders)
	mux.HandleFunc("POST /orders", orderHandler.handleCreateOrder)

	return mux
}
