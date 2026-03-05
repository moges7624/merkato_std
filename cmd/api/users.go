package main

import (
	"net/http"
	"strconv"

	"github.com/moges7624/merkato_std/internal/user"
)

type UserHandler struct {
	service user.Service
}

func NewUserHandler() *UserHandler {
	store := user.NewFileStore()

	return &UserHandler{
		service: *user.NewService(store),
	}
}

func (h *UserHandler) handleGetUsers(w http.ResponseWriter, _ *http.Request) {
	users, err := h.service.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"users": users})
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (h *UserHandler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "user id required", http.StatusBadRequest)
		return
	}
	if id < 1 {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUser(id)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	if user == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"user": user})
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (h *UserHandler) handleCreateUser(w http.ResponseWriter, _ *http.Request) {
	user, err := h.service.CreateUser()
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"user": user})
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (h *UserHandler) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "user id required", http.StatusBadRequest)
		return
	}

	if id < 1 {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	user, err := h.service.UpdateUser(id)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	if user == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"user": user})
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (h *UserHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "user id required", http.StatusBadRequest)
		return
	}

	if id < 1 {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	if err = h.service.DeleteUser(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{})
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
