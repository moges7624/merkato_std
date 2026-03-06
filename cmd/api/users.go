package main

import (
	"fmt"
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

func (h *UserHandler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetUsers()
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"users": users})
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}

func (h *UserHandler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		badRequestresponse(w, r, fmt.Errorf("user id required"))
		return
	}
	if id < 1 {
		badRequestresponse(w, r, fmt.Errorf("invalid user id"))
		return
	}

	user, err := h.service.GetUser(id)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	if user == nil {
		notFoundResponse(w, r)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"user": user})
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}

func (h *UserHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var input user.CreateUserParams

	err := readJSON(w, r, &input)
	if err != nil {
		badRequestresponse(w, r, err)
		return
	}

	user, err := h.service.CreateUser(&input)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"user": user})
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}

func (h *UserHandler) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		badRequestresponse(w, r, fmt.Errorf("user id required"))
		return
	}

	if id < 1 {
		badRequestresponse(w, r, fmt.Errorf("invalid user id"))
		return
	}

	var input user.UpateUserParams
	err = readJSON(w, r, &input)
	if err != nil {
		badRequestresponse(w, r, err)
		return
	}

	user, err := h.service.UpdateUser(id, input)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	if user == nil {
		notFoundResponse(w, r)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"user": user})
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}

func (h *UserHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		badRequestresponse(w, r, fmt.Errorf("user id required"))
		return
	}

	if id < 1 {
		badRequestresponse(w, r, fmt.Errorf("invalid user id"))
		return
	}

	if err = h.service.DeleteUser(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{})
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}
