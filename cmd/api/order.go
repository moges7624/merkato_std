package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/moges7624/merkato_std/internal/order"
	"github.com/moges7624/merkato_std/internal/product"
	"github.com/moges7624/merkato_std/internal/validator"
)

type orderHandler struct {
	service order.Service
	s       *APIServer
}

func NewOrderHandler(s *APIServer, service order.Service) *orderHandler {
	return &orderHandler{
		service: service,
		s:       s,
	}
}

func (h *orderHandler) handleGetOrders(
	w http.ResponseWriter,
	r *http.Request,
) {
	orders, err := h.service.GetOrders()
	if err != nil {
		h.s.serverErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"orders": orders})
	if err != nil {
		h.s.serverErrorResponse(w, r, err)
	}
}

func (h *orderHandler) handleGetOrderByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.s.badRequestresponse(w, r, fmt.Errorf("numeric order id required"))
		return
	}

	if id < 1 {
		h.s.badRequestresponse(w, r, fmt.Errorf("invalid order id"))
		return
	}

	o, err := h.service.GetOrderByID(int64(id))
	if err != nil {
		if errors.Is(err, order.ErrOrderNotFound) {
			h.s.notFoundResponse(w, r, err.Error())
		} else {
			h.s.serverErrorResponse(w, r, err)
		}
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"order": o})
	if err != nil {
		h.s.serverErrorResponse(w, r, err)
		return
	}
}

func (h *orderHandler) handleCreateOrder(
	w http.ResponseWriter,
	r *http.Request,
) {
	var input order.CreateOrderRequest

	err := readJSON(w, r, &input)
	if err != nil {
		h.s.badRequestresponse(w, r, err)
		return
	}

	v := validator.New()
	if order.ValidateCreateOrderRequest(v, &input); !v.Valid() {
		h.s.failedValidationResponse(w, r, v.Errors)
		return
	}

	order, err := h.service.CreateOrder(&input)
	if err != nil {
		if errors.Is(err, product.ErrInsufficientStock) {
			h.s.inventoryErrorResponse(w, r, err.Error())
		} else {
			h.s.badRequestresponse(w, r, err)
		}
		return
	}

	err = writeJSON(w, http.StatusCreated, envelope{"order": order})
	if err != nil {
		h.s.badRequestresponse(w, r, err)
	}
}
