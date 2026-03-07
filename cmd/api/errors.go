package main

import "net/http"

func errorResponse(
	w http.ResponseWriter,
	_ *http.Request,
	status int,
	message any,
) {
	msg := envelope{"error": message}

	err := writeJSON(w, status, msg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func serverErrorResponse(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	message := "Server encountered a problem and could not process your request"
	errorResponse(w, r, http.StatusInternalServerError, message)
}

func notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "Requested resource could not be found"
	errorResponse(w, r, http.StatusNotFound, message)
}

func badRequestresponse(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func failedValidationResponse(
	w http.ResponseWriter,
	r *http.Request,
	errors map[string]string,
) {
	errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
