package user

import (
	"errors"
	"sync"

	"github.com/brianvoe/gofakeit/v7"
)

type FileStore struct {
	mu    sync.RWMutex
	users map[int]User
	idSeq int
}

func NewFileStore() *FileStore {
	return &FileStore{
		users: map[int]User{
			1: {ID: 1, Name: "Arianna Banks", Email: gofakeit.Email()},
			2: {ID: 2, Name: "Nathanael Hale", Email: gofakeit.Email()},
		},
		idSeq: 2,
	}
}

func (s *FileStore) getUsers() (*[]User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}

	return &users, nil
}

func (s *FileStore) getUser(id int) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func (s *FileStore) getUserByEmail(email string) (*User, error) {
	return nil, errors.New("not impelemented")
	// s.mu.RLock()
	// defer s.mu.RUnlock()
	//
	// user, ok := s.users[id]
	// if !ok {
	// 	return nil, ErrUserNotFound
	// }
	// return &user, nil
}

func (s *FileStore) createUser(user *User) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.idSeq++
	user.ID = s.idSeq
	s.users[user.ID] = *user

	return user, nil
}

func (s *FileStore) updateUser(user User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.users[user.ID]
	if !ok {
		return ErrUserNotFound
	}

	s.users[user.ID] = user
	return nil
}

func (s *FileStore) deleteUser(id int) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	delete(s.users, id)

	return nil
}
