package main

import (
	"io"
	"net/http"

	"github.com/moges7624/merkato_std/internal/product"
)

func homeHandler(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "Hello there")
}

func (s *APIServer) NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/userHandler.llo", homeHandler)

	userHandler := NewUserHandler(s)
	mux.HandleFunc("GET /users", userHandler.handleGetUsers)
	mux.HandleFunc("POST /users", userHandler.handleCreateUser)
	mux.HandleFunc("GET /users/{id}", userHandler.handleGetUser)
	mux.HandleFunc("PATCH /users/{id}", userHandler.handleUpdateUser)
	mux.HandleFunc("DELETE /users/{id}", userHandler.handleDeleteUser)

	productFileStore := product.NewFileStore()
	productHandler := NewProductHandler(s, productFileStore)
	mux.HandleFunc("GET /products", productHandler.handleGetProducts)
	mux.HandleFunc("GET /products/{id}", productHandler.handleGetProduct)
	mux.HandleFunc("POST /products", productHandler.handleCreateProduct)

	return mux
}
