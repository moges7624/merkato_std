package user

type Store interface {
	getUsers() (*[]User, error)
	getUser(id int) (*User, error)
	createUser(user *User) (*User, error)
	updateUser(user User) error
	deleteUser(id int) error
}

type CreateUserParams struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpateUserParams struct {
	Name string `json:"name"`
}
