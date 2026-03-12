package main

import (
	"log/slog"
	"net/http"
)

type ErrorType string

const (
	ResourceNotFound    ErrorType = "resource_not_found"
	InvalidRequestError ErrorType = "invalid_request_error"
	InternalError       ErrorType = "api_error"
	InventoryError      ErrorType = "inventory_error"
)

type APIError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message,omitempty"`
	Details any       `json:"details,omitempty"`
}

func (s *APIServer) logError(apiError *APIError) {
	s.logger.Error(
		string(apiError.Message),
		slog.String("type", string(apiError.Type)),
	)
}

func (s *APIServer) errorResponse(
	w http.ResponseWriter,
	_ *http.Request,
	status int,
	apiError *APIError,
) {
	msg := envelope{"error": apiError}
	s.logError(apiError)

	err := writeJSON(w, status, msg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *APIServer) serverErrorResponse(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	message := "Server encountered a problem and could not process your request"
	APIErr := &APIError{
		Type:    InternalError,
		Message: message,
	}
	s.errorResponse(w, r, http.StatusInternalServerError, APIErr)
}

func (s *APIServer) notFoundResponse(
	w http.ResponseWriter,
	r *http.Request,
	msg string,
) {
	d := &APIError{
		Type:    ResourceNotFound,
		Message: msg,
	}
	s.errorResponse(w, r, http.StatusNotFound, d)
}

func (s *APIServer) badRequestresponse(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	APIErr := &APIError{
		Type:    InvalidRequestError,
		Message: err.Error(),
	}
	s.errorResponse(w, r, http.StatusUnprocessableEntity, APIErr)
}

func (s *APIServer) failedValidationResponse(
	w http.ResponseWriter,
	r *http.Request,
	errors map[string]string,
) {
	APIErr := &APIError{
		Type:    InvalidRequestError,
		Message: "validation failed",
		Details: errors,
	}
	s.errorResponse(w, r, http.StatusUnprocessableEntity, APIErr)
}

func (s *APIServer) inventoryErrorResponse(
	w http.ResponseWriter,
	r *http.Request,
	message string,
) {
	APIErr := &APIError{
		Type:    InventoryError,
		Message: message,
	}
	s.errorResponse(w, r, http.StatusUnprocessableEntity, APIErr)
}
