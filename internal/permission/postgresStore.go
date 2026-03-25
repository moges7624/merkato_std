package permission

import "database/sql"

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{
		db: db,
	}
}

func (ps *PostgresStore) getAll() ([]Permission, error) {
	query := `
	SELECT code from permissions
	`

	rows, err := ps.db.Query(query)
	if err != nil {
		return nil, err
	}

	var permissions []Permission
	for rows.Next() {
		var p Permission

		err = rows.Scan(&p.Code)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (ps *PostgresStore) getAllForUser(userID int64) ([]Permission, error) {
	query := `
	SELECT permissions.code from permissions 
	INNER JOIN users_permissions
	ON permissions.id = users_permissions.permission_id
	WHERE users_permissions.user_id = $1
	`

	rows, err := ps.db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	var permissions []Permission
	for rows.Next() {
		var p Permission

		err = rows.Scan(&p.Code)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (ps *PostgresStore) addForUser(userID int64, pIds []int64) error {
	query := `
	INSERT INTO users_permissions (user_id, permission_id)
	SELECT $1, id
	FROM permissions
	WHERE id = ANY($2::int[])
	ON CONFLICT DO NOTHING`

	_, err := ps.db.Exec(query, userID, pIds)
	return err
}
