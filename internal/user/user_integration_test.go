//go:build integration

package user

import (
	"slices"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/moges7624/merkato_std/internal/assert"
	"github.com/moges7624/merkato_std/internal/utils"
)

func newUserService(t *testing.T) *Service {
	db := utils.NewTestDB(t)
	return NewService(NewPostgresStore(db))
}

func TestUserService_GetUsers_Empty(t *testing.T) {
	if testing.Short() {
		t.Skip("UserService: skipping integration test")
	}

	userSvc := newUserService(t)

	t.Run("return empty slice", func(t *testing.T) {
		users, err := userSvc.GetUsers()
		if err != nil {
			t.Fatal(err)
		}

		if !slices.Equal(users, []*User{}) {
			t.Errorf("got: %v, expected: %v\n", users, []*User{})
		}
	})
}

func TestUserService_CreateUser(t *testing.T) {
	if testing.Short() {
		t.Skip("UserService: skipping integration test")
	}

	userSvc := newUserService(t)

	name := gofakeit.Name()
	email := gofakeit.Email()
	password := gofakeit.Password(true, true, false, false, false, 8)

	createUserIput := CreateUserParams{
		Name:              name,
		Email:             email,
		PlainTextPassword: password,
	}

	var user *User

	t.Run("success", func(t *testing.T) {
		u, err := userSvc.CreateUser(&createUserIput)

		assert.Null(t, err)
		assert.NotNull(t, u)
		assert.Equal(t, u.ID, 1)
		assert.Equal(t, u.Name, name)

		if *u.Password.plaintext == string(u.Password.hash) {
			t.Errorf("password not hashed properly")
		}

		user = u
	})

	t.Run("duplicate email", func(t *testing.T) {
		u, err := userSvc.CreateUser(&createUserIput)

		assert.NotNull(t, err)
		assert.Null(t, u)
		assert.Equal(t, err, ErrUserAlreadyExists)
	})

	t.Run("password does not match", func(t *testing.T) {
		ok, err := user.Password.PasswordMatches("randompass")

		assert.Null(t, err)
		assert.Equal(t, ok, false)
	})

	t.Run("password matches", func(t *testing.T) {
		ok, err := user.Password.PasswordMatches(password)

		assert.Null(t, err)
		assert.Equal(t, ok, true)
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	if testing.Short() {
		t.Skip("UserService: skipping integration test")
	}

	userSvc := newUserService(t)

	updateUserParams := UpateUserParams{
		Name: "john",
	}

	t.Run("non-existing user", func(t *testing.T) {
		u, err := userSvc.UpdateUser(1, updateUserParams)

		assert.NotNull(t, err)
		assert.Equal(t, err, ErrUserNotFound)
		assert.Null(t, u)
	})

	t.Run("existing user", func(t *testing.T) {
		name := gofakeit.Name()
		email := gofakeit.Email()
		createUserIput := CreateUserParams{
			Name:              name,
			Email:             email,
			PlainTextPassword: gofakeit.Password(false, false, false, false, false, 8),
		}

		_, err := userSvc.CreateUser(&createUserIput)
		if err != nil {
			t.Fatal(err)
		}

		u, err := userSvc.UpdateUser(1, updateUserParams)

		assert.Null(t, err)
		assert.NotNull(t, u)
		assert.Equal(t, u.Name, "john")
	})
}
