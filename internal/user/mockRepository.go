package user

import "github.com/stretchr/testify/mock"

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) getUsers() ([]*User, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*User), args.Error(1)
}

func (m *MockRepository) getUser(id int) (*User, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) getUserByEmail(email string) (*User, error) {
	args := m.Called(email)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) createUser(u *User) (*User, error) {
	args := m.Called(u)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) updateUser(u User) error {
	args := m.Called(u)

	return args.Error(1)
}

func (m *MockRepository) deleteUser(id int) error {
	args := m.Called(id)

	return args.Error(1)
}
