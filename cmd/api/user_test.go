//go:build integration

package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/moges7624/merkato_std/internal/assert"
	"github.com/moges7624/merkato_std/internal/user"
	"github.com/moges7624/merkato_std/internal/utils"
)

func newUserHandler(t *testing.T) *UserHandler {
	db := utils.NewTestDB(t)
	api := &APIServer{
		logger: slog.New(slog.DiscardHandler),
	}

	userRepo := user.NewPostgresStore(db)
	userSvc := user.NewService(userRepo)
	return NewUserHandler(api, *userSvc)
}

func TestUserHandler_getUsers(t *testing.T) {
	userHandler := newUserHandler(t)

	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	userHandler.handleGetUsers(w, req)

	res := w.Result()

	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func TestUserHandler_CreateUser(t *testing.T) {
	UserHandler := newUserHandler(t)

	t.Run("given incomplete input, it should return 422", func(t *testing.T) {
		body := `{"name": "john", "email":"john@mail.com"}`
		req, err := http.NewRequest(http.MethodGet, "/", strings.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		UserHandler.handleCreateUser(w, req)

		res := w.Result()

		assert.Equal(t, res.StatusCode, http.StatusUnprocessableEntity)
	})

	t.Run("given invalid email, it should return 422", func(t *testing.T) {
		body := `{"name": "john", "email":"johnmail.com", "password": "12345678"}`
		req, err := http.NewRequest(http.MethodGet, "/", strings.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		UserHandler.handleCreateUser(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, res.StatusCode, http.StatusUnprocessableEntity)

		var apiRes struct {
			Error APIError
		}

		err = json.NewDecoder(res.Body).Decode(&apiRes)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, apiRes.Error.Message, "validation failed")
	})

	t.Run("given valid input, it should return 200", func(t *testing.T) {
		body := `{"name": "john", "email":"john@mail.com", "password":"12345678"}`
		req, err := http.NewRequest(http.MethodGet, "/", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		UserHandler.handleCreateUser(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, res.StatusCode, http.StatusOK)

		var apiRes struct {
			User user.User
		}

		err = json.NewDecoder(res.Body).Decode(&apiRes)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, apiRes.User.ID, 2)
		assert.Equal(t, apiRes.User.Email, "john@mail.com")
	})

	t.Run("given a user exists, it should return 422", func(t *testing.T) {
		body := `{"name": "john", "email":"john@mail.com", "password":"12345678"}`
		req, err := http.NewRequest(http.MethodGet, "/", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		UserHandler.handleCreateUser(w, req)
		UserHandler.handleCreateUser(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, res.StatusCode, http.StatusUnprocessableEntity)

		var apiRes struct {
			Error APIError
		}

		err = json.NewDecoder(res.Body).Decode(&apiRes)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, apiRes.Error.Message, "user with given info already exists")
	})
}

func TestUserHandler_UpdateUser(t *testing.T) {
	UserHandler := newUserHandler(t)

	t.Run("given invalid user id, it should return 422", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("POST /users/{id}", UserHandler.handleUpdateUser)

		req := httptest.NewRequest("POST", "/users/1a", nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode == http.StatusMethodNotAllowed {
			t.Fatal("method not allowed")
		}

		var apiRes struct {
			Error APIError
		}

		err := json.NewDecoder(res.Body).Decode(&apiRes)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, apiRes.Error.Message, "invalid user id")
	})

	t.Run("given non existing user id, it should return 404", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("PATCH /users/{id}", UserHandler.handleUpdateUser)

		reqBody := `{"name": "james"}`
		req := httptest.NewRequest("PATCH", "/users/2", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode == http.StatusMethodNotAllowed {
			t.Fatal("method not allowed")
		}

		var apiRes struct {
			Error APIError
		}

		err := json.NewDecoder(res.Body).Decode(&apiRes)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, apiRes.Error.Message, user.ErrUserNotFound.Error())
		assert.Equal(t, apiRes.Error.Type, ResourceNotFound)
	})

	t.Run("given existing user id, it should return 200", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("PATCH /users/{id}", UserHandler.handleUpdateUser)

		reqBody := `{"name": "james"}`
		req := httptest.NewRequest("PATCH", "/users/1", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode == http.StatusMethodNotAllowed {
			t.Fatal("method not allowed")
		}

		assert.Equal(t, res.StatusCode, http.StatusOK)

		var apiRes struct {
			User user.User
		}

		err := json.NewDecoder(res.Body).Decode(&apiRes)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, apiRes.User.Name, "james")
	})
}
