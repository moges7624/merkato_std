package main

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/moges7624/merkato_std/internal/auth"
)

type contextKey string

const UserClaimsKey contextKey = "user_claims"

func (s *APIServer) AuthRequired(
	jwtSvc *auth.JWTService,
	next http.HandlerFunc,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			s.failedAuthenticationResponse(w, r, "missing authentication token")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			s.failedAuthenticationResponse(w, r, "invalid authentication token")
			return
		}

		claims, err := jwtSvc.ValidateAccessToken(parts[1])
		if err != nil {
			if errors.Is(err, auth.ErrExpiredToken) {
				s.failedAuthenticationResponse(w, r, err.Error())
			} else {
				s.failedAuthenticationResponse(w, r, "authentication failed")
			}
			return
		}

		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
