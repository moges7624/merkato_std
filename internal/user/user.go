package user

import (
	"errors"
	"regexp"

	"github.com/moges7624/merkato_std/internal/validator"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type User struct {
	ID    int
	Name  string
	Email string
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 72 bytes long")
	v.Check(user.Email != "", "email", "must be provided")
	v.Check(
		validator.Matches(user.Email, regexp.MustCompile(validator.EmailRX)),
		"email",
		"must be a valid email address",
	)
}
