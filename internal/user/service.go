package user

import (
	"fmt"
)

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) GetUsers() (*[]User, error) {
	return s.store.getUsers()
}

func (s *Service) GetUser(id int) (*User, error) {
	user, err := s.store.getUser(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) CreateUser(tmpUser *CreateUserParams) (*User, error) {
	user := &User{
		Name:  tmpUser.Name,
		Email: tmpUser.Email,
	}

	user, err := s.store.createUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) UpdateUser(
	id int,
	updateInput UpateUserParams,
) (*User, error) {
	user, err := s.store.getUser(id)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if updateInput.Name != "" {
		user.Name = updateInput.Name
	}

	if err := s.store.updateUser(*user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) DeleteUser(id int) error {
	_, err := s.store.getUser(id)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if err := s.store.deleteUser(id); err != nil {
		return err
	}

	return nil
}
