package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/moges7624/merkato_std/internal/filter"
	"github.com/moges7624/merkato_std/internal/product"
	"github.com/moges7624/merkato_std/internal/validator"
)

type ProductHandler struct {
	service product.Service
	s       *APIServer
}

func NewProductHandler(
	s *APIServer,
	service product.Service,
) *ProductHandler {
	return &ProductHandler{
		service: service,
		s:       s,
	}
}

func (h *ProductHandler) handleGetProducts(
	w http.ResponseWriter,
	r *http.Request,
) {
	var input struct {
		Name string
		filter.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Name = h.s.readString(qs, "name", "")
	input.Page = h.s.readInt(qs, "page", 1, v)
	input.PageSize = h.s.readInt(qs, "page_size", 20, v)
	input.Sort = h.s.readString(qs, "sort", "id")

	input.SortSafelist = []string{
		"name", "-name", "id", "-id",
	}

	if filter.ValidateFilters(v, input.Filters); !v.Valid() {
		h.s.failedValidationResponse(w, r, v.Errors)
		return
	}

	prods, metadata, err := h.service.GetProducts(&product.ProductFilters{
		Name:    input.Name,
		Filters: input.Filters,
	})
	if err != nil {
		h.s.serverErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"products": prods, "metadata": metadata})
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
