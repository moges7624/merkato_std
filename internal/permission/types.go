package permission

type Store interface {
	getAll() ([]Permission, error)
	getAllForUser(id int64) ([]Permission, error)
	addForUser(userID int64, pIDs []int64) error
}

type AddPermissionForUserRequest struct {
	UserID      int64   `json:"user_id"`
	Permissions []int64 `json:"permissions"`
}
