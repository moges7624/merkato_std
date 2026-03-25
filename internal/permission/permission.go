package permission

import "errors"

var ErrPermissionNotFound = errors.New("permission not found")

type Permission struct {
	Code string `json:"code"`
}
