package postgres

import (
	"context"
	"database/sql"

	"github.com/DeepayanMallick/base-grpc/usermgm/storage"
)

const saltLen = 32

const insertUser = `
INSERT INTO users (
	first_name,
	last_name,
	username,
	email,
	password,
	status
) VALUES (
	:first_name,
	:last_name,  
	:username,
	:email,
	:password,
	:status
) RETURNING
	id
`

func (s *Storage) CreateUser(ctx context.Context, user storage.User) (string, error) {
	stmt, err := s.db.PrepareNamed(insertUser)
	if err != nil {
		return "", err
	}
	var id string
	if err := stmt.Get(&id, user); err != nil {
		return "", err
	}

	return id, nil
}

const getUserByEmail = `
SELECT
	*
FROM
	users
WHERE
	email = :email
AND
	deleted_at IS NULL
`

func (s Storage) GetUserByEmail(ctx context.Context, email string) (*storage.User, error) {
	var u storage.User
	stmt, err := s.db.PrepareNamed(getUserByEmail)
	if err != nil {
		return nil, err
	}

	arg := map[string]interface{}{
		"email": email,
	}
	if err := stmt.Get(&u, arg); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}

	return &u, nil
}

const getUserByID = `
SELECT
	*
FROM
	users
WHERE
	id = $1
AND
	deleted_at IS NULL
`

func (s Storage) GetUserByID(ctx context.Context, id string) (*storage.User, error) {
	var u storage.User
	err := s.db.Get(&u, getUserByID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}

	return &u, nil
}

const deleteUser = `DELETE FROM users where id = $1`

func (s Storage) deleteUserPermanently(ctx context.Context, id string) error {
	row, err := s.db.Exec(deleteUser, id)
	if err != nil {
		s.logger.WithError(err)
		return err
	}

	rowCount, err := row.RowsAffected()
	if err != nil {
		s.logger.WithError(err)
		return err
	}

	if rowCount <= 0 {
		s.logger.Error("Unable to delete the targeted user")
		return storage.NotFound
	}

	return nil
}

const deleteAllUsers = `DELETE FROM users`

func (s Storage) deleteUsersPermanently(ctx context.Context) error {
	row, err := s.db.ExecContext(ctx, deleteAllUsers)
	if err != nil {
		s.logger.WithError(err)
		return err
	}

	rowCount, err := row.RowsAffected()
	if err != nil {
		s.logger.WithError(err)
		return err
	}

	if rowCount <= 0 {
		s.logger.Error("Unable to delete users")
		return storage.NotFound
	}

	return nil
}
