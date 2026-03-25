package main

import (
	"io"
	"net/http"

	"github.com/moges7624/merkato_std/internal/auth"
	"github.com/moges7624/merkato_std/internal/order"
	"github.com/moges7624/merkato_std/internal/permission"
	"github.com/moges7624/merkato_std/internal/product"
	"github.com/moges7624/merkato_std/internal/user"
)

func homeHandler(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "Hello there")
}

func (s *APIServer) NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/home", homeHandler)

	userPgStore := user.NewPostgresStore(s.DB)
	userService := user.NewService(userPgStore)

	// TODO: read the secretes from env
	authService := auth.NewJWTService("sdsdfdJljjadfef", "sdsdfdJljjadfef")
	authHandler := NewAuthHandler(s, *authService, *userService)
	mux.HandleFunc("POST /auth/login", authHandler.handleLogin)

	permissionRepo := permission.NewPostgresStore(s.DB)
	permissionService := permission.NewService(permissionRepo)
	permissionHandler := NewPermissionHandler(s, *permissionService)
	mux.HandleFunc("GET /permissions", permissionHandler.handleGetAll)
	mux.HandleFunc("POST /permissions", permissionHandler.handleAddForUser)

	// userFileStore := user.NewFileStore()
	userHandler := NewUserHandler(s, *userService)
	mux.HandleFunc("GET /users",
		s.AuthRequired(authService, userHandler.handleGetUsers))
	mux.HandleFunc("POST /users", userHandler.handleCreateUser)
	mux.HandleFunc("GET /users/{id}",
		s.AuthRequired(authService, userHandler.handleGetUser))
	mux.HandleFunc("PATCH /users/{id}",
		s.AuthRequired(authService, userHandler.handleUpdateUser))
	mux.HandleFunc("DELETE /users/{id}",
		s.AuthRequired(authService, userHandler.handleDeleteUser))

	// productFileStore := product.NewFileStore()
	productPostgresStore := product.NewPostgresStore(s.DB)
	productService := product.NewService(productPostgresStore)
	productHandler := NewProductHandler(s, *productService)
	mux.HandleFunc("GET /products",
		s.AuthRequired(authService, productHandler.handleGetProducts))
	mux.HandleFunc("GET /products/{id}",
		s.AuthRequired(authService, productHandler.handleGetProduct))
	mux.HandleFunc("POST /products",
		s.AuthRequired(authService, productHandler.handleCreateProduct))

	// orderFileStore := order.NewFileStore()
	orderPostgresStore := order.NewPostgresStore(s.DB)
	orderService := order.NewService(
		orderPostgresStore,
		*productService,
		*userService)
	orderHandler := NewOrderHandler(s, *orderService)
	mux.HandleFunc("GET /orders",
		s.AuthRequired(authService, orderHandler.handleGetOrders))
	mux.HandleFunc("POST /orders",
		s.AuthRequired(authService, orderHandler.handleCreateOrder))
	mux.HandleFunc("GET /orders/{id}",
		s.AuthRequired(authService, orderHandler.handleGetOrderByID))

	return mux
}
