package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/moges7624/merkato_std/internal/auth"
	"github.com/moges7624/merkato_std/internal/user"
	"github.com/moges7624/merkato_std/internal/validator"
)

type authHandler struct {
	authService auth.JWTService
	userService user.Service
	s           *APIServer
}

func NewAuthHandler(s *APIServer,
	jwtService auth.JWTService,
	userService user.Service,
) *authHandler {
	return &authHandler{
		authService: jwtService,
		userService: userService,
		s:           s,
	}
}

func (h *authHandler) handleLogin(
	w http.ResponseWriter,
	r *http.Request,
) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := readJSON(w, r, &input); err != nil {
		h.s.badRequestresponse(w, r, err)
		return
	}

	v := validator.New()

	user.ValidatePasswordPlaintext(v, input.Password)
	user.ValidateEmail(v, input.Email)

	if valid := v.Valid(); !valid {
		h.s.failedValidationResponse(w, r, v.Errors)
		return
	}
	u, err := h.userService.GetUserByEmail(input.Email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			h.s.failedAuthenticationResponse(w, r, "invalid credentials")
		} else {
			h.s.serverErrorResponse(w, r, err)
		}
		return
	}

	match, err := u.Password.PasswordMatches(input.Password)
	if err != nil {
		h.s.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		h.s.failedAuthenticationResponse(w, r, "invalid credentials")
		return
	}

	accessToken, err := h.authService.GenerateAccessToken(strconv.Itoa(u.ID))
	if err != nil {
		h.s.serverErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"access_token": accessToken})
	if err != nil {
		h.s.serverErrorResponse(w, r, err)
		return
	}
}
