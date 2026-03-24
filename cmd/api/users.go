package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	usr "github.com/moges7624/merkato_std/internal/user"
	"github.com/moges7624/merkato_std/internal/validator"
)

type UserHandler struct {
	service usr.Service
	s       *APIServer
}

func NewUserHandler(s *APIServer, userService usr.Service) *UserHandler {
	return &UserHandler{
		service: userService,
		s:       s,
	}
}

func (h *UserHandler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetUsers()
	if err != nil {
		h.s.serverErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"users": users})
	if err != nil {
		h.s.serverErrorResponse(w, r, err)
	}
}

func (h *UserHandler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.s.badRequestresponse(w, r, fmt.Errorf("user id required"))
		return
	}

	if id < 1 {
		h.s.badRequestresponse(w, r, fmt.Errorf("invalid user id"))
		return
	}

	user, err := h.service.GetUser(id)
	if err != nil {
		if errors.Is(err, usr.ErrUserNotFound) {
			h.s.notFoundResponse(w, r, err.Error())
		} else {
			h.s.serverErrorResponse(w, r, err)
		}
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"user": user})
	if err != nil {
		h.s.serverErrorResponse(w, r, err)
	}
}

func (h *UserHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var input usr.CreateUserParams

	err := readJSON(w, r, &input)
	if err != nil {
		h.s.badRequestresponse(w, r, err)
		return
	}

	v := validator.New()

	if input.Validate(v); !v.Valid() {
		h.s.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := h.service.CreateUser(&input)
	if err != nil {
		if errors.Is(err, usr.ErrUserAlreadyExists) {
			h.s.badRequestresponse(w, r, fmt.Errorf("user with given info already exists"))
		} else {
			h.s.serverErrorResponse(w, r, err)
		}
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"user": user})
	if err != nil {
		h.s.serverErrorResponse(w, r, err)
	}
}

func (h *UserHandler) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		h.s.badRequestresponse(w, r, fmt.Errorf("invalid user id"))
		return
	}

	var input usr.UpateUserParams
	err = readJSON(w, r, &input)
	if err != nil {
		h.s.badRequestresponse(w, r, err)
		return
	}

	v := validator.New()

	v.Check(input.Name != "", "name", "must be provided")

	if !v.Valid() {
		h.s.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := h.service.UpdateUser(id, input)
	if err != nil {
		if errors.Is(err, usr.ErrUserNotFound) {
			h.s.notFoundResponse(w, r, err.Error())
		} else {
			h.s.serverErrorResponse(w, r, err)
		}
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"user": user})
	if err != nil {
		h.s.serverErrorResponse(w, r, err)
	}
}

func (h *UserHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.s.badRequestresponse(w, r, fmt.Errorf("user id required"))
		return
	}

	if id < 1 {
		h.s.badRequestresponse(w, r, fmt.Errorf("invalid user id"))
		return
	}

	if err = h.service.DeleteUser(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{})
	if err != nil {
		h.s.serverErrorResponse(w, r, err)
	}
}
