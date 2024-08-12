package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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

	// Hash the plain-text password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	// Insert user credentials, including hashed password, into database
	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))

	// Check if the error is a MySQL error
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

// Authenticate a user based on email and password; if the user exists, return ID
func (m *UserModel) Authenticate(email, password string) (int, error) {
	// Return ID and hashed password associated with given email
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, nil
		}
	}

	// Check if hashed password and plain-text password match
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// Else, password is correct; return user ID
	return id, nil

}

// Checks if user with specific ID exists
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
