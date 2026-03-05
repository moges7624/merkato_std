package main

import (
	"io"
	"net/http"
)

func homeHandler(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "Hello there")
}

func (s *APIServer) NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/userHandler.llo", homeHandler)

	userHandler := NewUserHandler()
	mux.HandleFunc("GET /users", userHandler.handleGetUsers)
	mux.HandleFunc("POST /users", userHandler.handleCreateUser)
	mux.HandleFunc("GET /users/{id}", userHandler.handleGetUser)
	mux.HandleFunc("PATCH /users/{id}", userHandler.handleUpdateUser)
	mux.HandleFunc("DELETE /users/{id}", userHandler.handleDeleteUser)

	return mux
}
