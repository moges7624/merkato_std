package user

import (
	"regexp"

	"github.com/moges7624/merkato_std/internal/validator"
)

type Store interface {
	getUsers() ([]*User, error)
	getUser(id int) (*User, error)
	getUserByEmail(email string) (*User, error)
	createUser(user *User) (*User, error)
	updateUser(user User) error
	deleteUser(id int) error
}

type CreateUserParams struct {
	Name              string `json:"name"`
	Email             string `json:"email"`
	PlainTextPassword string `json:"password"`
}

type UpateUserParams struct {
	Name string `json:"name"`
}

func (cup *CreateUserParams) Validate(v *validator.Validator) {
	v.Check(cup.Name != "", "name", "must be provided")
	v.Check(len(cup.Name) <= 500, "name", "must not be more than 72 bytes long")
	v.Check(cup.Email != "", "email", "must be provided")
	v.Check(cup.PlainTextPassword != "", "password", "must be provided")

	if cup.PlainTextPassword != "" {
		ValidatePasswordPlaintext(v, cup.PlainTextPassword)
	}

	v.Check(
		validator.Matches(cup.Email, regexp.MustCompile(validator.EmailRX)),
		"email",
		"must be a valid email address",
	)
}
