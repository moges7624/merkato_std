package user

type Store interface {
	getUsers() (*[]User, error)
	getUser(id int) (*User, error)
	createUser() (*User, error)
	updateUser(user User) error
	deleteUser(id int) error
}
