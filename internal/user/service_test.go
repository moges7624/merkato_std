package user

import (
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/moges7624/merkato_std/internal/assert"
)

func TestGetUsers(t *testing.T) {
	fileStore := NewFileStore()
	userSevice := NewService(fileStore)

	t.Run("return all users", func(t *testing.T) {
		users, err := userSevice.GetUsers()
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, len(users), 2)
	})
}

func TestGetUser(t *testing.T) {
	fileStore := NewFileStore()
	userSevice := NewService(fileStore)

	t.Run("Given existing user id, it should return a user", func(t *testing.T) {
		user, err := userSevice.GetUser(2)
		if err != nil {
			t.Fatal(err)
		}

		assert.NotNull(t, user)
	})

	t.Run("Given non existing user id, it should throw not found error", func(t *testing.T) {
		user, err := userSevice.GetUser(3)
		if err == nil {
			t.Fatal("error shouldn't be nil")
		}

		assert.Null(t, user)

		if !errors.Is(err, ErrUserNotFound) {
			t.Errorf("got %s expected %s", err.Error(), ErrUserNotFound.Error())
		}
	})
}

func TestCreateUser(t *testing.T) {
	t.Run("given a valid user struct, it should create a user", func(t *testing.T) {
		name := gofakeit.Name()
		email := gofakeit.Email()
		createUserParam := &CreateUserParams{
			Name:  name,
			Email: email,
		}

		fileStore := NewFileStore()
		userSevice := NewService(fileStore)

		user, err := userSevice.CreateUser(createUserParam)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, user.Email, email)
		assert.Equal(t, user.Name, name)
	})
}
