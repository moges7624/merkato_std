package permission

import (
	"errors"
)

var ErrPermissionNotFound = errors.New("permission not found")

type Permission struct {
	Code string `json:"code"`
}

func Includes(permissions []Permission, code string) bool {
	for i := range permissions {
		if permissions[i].Code == code {
			return true
		}
	}

	return false
}
