package main

import (
	"net/http"

	"github.com/moges7624/merkato_std/internal/permission"
)

type PermissinHandler struct {
	service permission.Service
	s       *APIServer
}

func NewPermissionHandler(s *APIServer, svc permission.Service) *PermissinHandler {
	return &PermissinHandler{
		service: svc,
	}
}

func (h *PermissinHandler) handleGetAll(
	w http.ResponseWriter,
	r *http.Request,
) {
	permissions, err := h.service.GetAll()
	if err != nil {
		h.s.serverErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"permissions": permissions})
	if err != nil {
		h.s.serverErrorResponse(w, r, err)
	}
}

func (h *PermissinHandler) handleAddForUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	var input permission.AddPermissionForUserRequest
	err := readJSON(w, r, &input)
	if err != nil {
		h.s.serverErrorResponse(w, r, err)
		return
	}

	// TODO: validate user input

	err = h.service.AddForUser(input)
	if err != nil {
		h.s.serverErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, nil)
	if err != nil {
		h.s.serverErrorResponse(w, r, err)
	}
}
