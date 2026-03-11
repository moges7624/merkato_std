package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/moges7624/merkato_std/internal/product"
	"github.com/moges7624/merkato_std/internal/validator"
)

type ProductHandler struct {
	service product.Service
	s       *APIServer
}

func NewProductHandler(s *APIServer, store product.Store) *ProductHandler {
	return &ProductHandler{
		service: *product.NewService(store),
		s:       s,
	}
}

func (h *ProductHandler) handleGetProducts(
	w http.ResponseWriter,
	r *http.Request,
) {
	prods, err := h.service.GetProducts()
	if err != nil {
		h.s.serverErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"products": prods})
	if err != nil {
		h.s.serverErrorResponse(w, r, err)
		return
	}
}

func (h *ProductHandler) handleGetProduct(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.s.badRequestresponse(w, r, fmt.Errorf("product id param must be an int"))
		return
	}

	if id < 1 {
		h.s.badRequestresponse(w, r, fmt.Errorf("invalid product id"))
		return
	}

	prod, err := h.service.GetProduct(id)
	if err != nil {
		if errors.Is(err, product.ErrProductNotFound) {
			h.s.notFoundResponse(w, r, err.Error())
		} else {
			h.s.serverErrorResponse(w, r, err)
		}
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"product": prod})
	if err != nil {
		h.s.serverErrorResponse(w, r, err)
		return
	}
}

func (h *ProductHandler) handleCreateProduct(
	w http.ResponseWriter,
	r *http.Request,
) {
	var input product.CreateProductRequest

	err := readJSON(w, r, &input)
	if err != nil {
		h.s.badRequestresponse(w, r, err)
		return
	}

	v := validator.New()
	if product.ValidateCreateProductParams(v, input); !v.Valid() {
		h.s.failedValidationResponse(w, r, v.Errors)
		return
	}

	prod, err := h.service.CreateProducts(&input)
	if err != nil {
		h.s.serverErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusCreated, envelope{"product": prod})
	if err != nil {
		h.s.serverErrorResponse(w, r, err)
	}
}
