package storage

import (
	"database/sql"
	"errors"
	"time"
)

const (
	Pass     string = "PASS"
	TOTP     string = "TOTP"
	PINCode  string = "CODE"
	SMS      string = "SMS"
	Recovery string = "RECOVERY"
	EMail    string = "EMAIL"
)

var (
	// NotFound is returned when the requested resource does not exist.
	NotFound = errors.New("not found")
	// Conflict is returned when trying to create the same resource twice.
	Conflict = errors.New("conflict")
	// UsernameExists is returned when the username already exists in storage.
	UsernameExists = errors.New("username already exists")
	// EmailExists is returned when signup email already exists in storage.
	EmailExists = errors.New("email already exists")
	// InvCodeExists is returned when invitation code already exists in storage.
	InvCodeExists = errors.New("invitation code already exists")
)

type (
	User struct {
		ID        string         `db:"id"`
		FirstName string         `db:"first_name"`
		LastName  string         `db:"last_name"`
		Username  string         `db:"username"`
		Email     string         `db:"email"`
		Password  string         `db:"password"`
		Status    int32          `db:"status"`
		CreatedAt time.Time      `db:"created_at"`
		UpdatedAt time.Time      `db:"updated_at"`
		DeletedAt sql.NullTime   `db:"deleted_at,omitempty"`
		DeletedBy sql.NullString `db:"deleted_by,omitempty"`
	}

	FilterUser struct {
		SearchTerm   string
		Limit        int32
		Offset       int32
		SortBy       string
		SortByColumn string
		Status       []string
	}
)
