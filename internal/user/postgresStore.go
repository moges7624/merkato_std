package user

import (
	"database/sql"
	"errors"
	"strings"
)

type PostgresStore struct {
	DB *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{
		DB: db,
	}
}

func (ps *PostgresStore) getUsers() (*[]User, error) {
	query := `SELECT id, name, email 
	FROM users`

	rows, err := ps.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var user User

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return &users, nil
}

func (ps *PostgresStore) getUser(id int) (*User, error) {
	query := `SELECT id, name, email 
	FROM users 
	WHERE id = $1`
	var user User

	err := ps.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrUserNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (ps *PostgresStore) createUser(user *User) (*User, error) {
	query := `INSERT INTO users (name, email)
	VALUES ($1, $2)
	RETURNING id`

	err := ps.DB.QueryRow(query, user.Name, user.Email).Scan(
		&user.ID)
	if err != nil {
		switch {
		case strings.HasPrefix(
			err.Error(),
			`pq: duplicate key value violates unique constraint "users_email_key"`,
		):
			return nil, ErrUserAlreadyExists
		default:
			return nil, err
		}
	}

	return user, nil
}

func (ps *PostgresStore) updateUser(user User) error {
	query := `UPDATE users
	SET name = $1, email = $2
	WHERE id = $3`

	_, err := ps.DB.Exec(query, user.Name, user.Email, user.ID)
	if err != nil {
		switch {
		case strings.HasPrefix(
			err.Error(),
			`pq: duplicate key value violates unique constraint "users_email_key"`,
		):
			return ErrUserAlreadyExists
		case errors.Is(err, sql.ErrNoRows):
			return ErrUserNotFound
		default:
			return err
		}
	}
	return nil
}

func (ps *PostgresStore) deleteUser(id int) error {
	return nil
}
