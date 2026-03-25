package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/moges7624/merkato_std/internal/auth"
	"github.com/moges7624/merkato_std/internal/permission"
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
		userID, ok := ctx.Value("userID").(string)
		if !ok {
			panic("missing user value in request context")
		}
		fmt.Println("authenticate: user id from context is:", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (s *APIServer) requirePermission(
	code string,
	next http.HandlerFunc,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("userID").(string)
		if !ok {
			panic("missing user value in request context")
		}

		uID, err := strconv.Atoi(userID)
		if err != nil {
			panic("invalid user id")
		}
		permissionRepo := permission.NewPostgresStore(s.DB)
		permissionService := permission.NewService(permissionRepo)
		permissions, err := permissionService.GetAllForUser(int64(uID))
		if err != nil {
			s.serverErrorResponse(w, r, err)
			return
		}

		if !permission.Includes(permissions, code) {
			s.failedAuthorizationResponse(w, r, "authorization failed")
			return
		}

		next.ServeHTTP(w, r)
	}

	// return s.AuthRequired(fn)
}
