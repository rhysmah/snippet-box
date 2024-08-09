package models

import (
	"database/sql"
	"time"
)

// User struct that exactly mirrors the database
// representation of a user.
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// Authenticate a user based on email and password; if the user exists, return ID
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Checks if user with specific ID exists
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
