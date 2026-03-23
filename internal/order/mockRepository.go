package order

import "github.com/stretchr/testify/mock"

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) getAll() ([]*Order, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*Order), args.Error(1)
}

func (m *MockRepo) getByID(id int64) (*Order, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Order), args.Error(1)
}

func (m *MockRepo) insert(o *Order) error {
	args := m.Called(o)

	return args.Error(0)
}
