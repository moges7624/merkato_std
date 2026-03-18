package user

import (
	"errors"
	"regexp"

	"github.com/moges7624/merkato_std/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type password struct {
	plaintext *string
	hash      []byte
}

type User struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password password `json:"-"`
}

func (u *User) setPassword(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	u.Password = password{
		plaintext: &plaintextPassword,
		hash:      hash,
	}
	return nil
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 72 bytes long")
	v.Check(user.Email != "", "email", "must be provided")

	// if user.Password.plaintext != nil {
	// 	ValidatePasswordPlaintext(v, *user.Password.plaintext)
	// }

	// if user.Password.hash == nil {
	// 	panic("missing password hash for user")
	// }

	v.Check(
		validator.Matches(user.Email, regexp.MustCompile(validator.EmailRX)),
		"email",
		"must be a valid email address",
	)
}
