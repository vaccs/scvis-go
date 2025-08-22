package models

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModelInterface interface {
	Insert(username, email, password string, administrator bool) error
	Authenticate(emailOrUsername, password string) (int, error)
	Exists(id int) (bool, error)
	IsAdministrator(id int) (bool, error)
	ChangePassword(id int, currentPassword, newPassword string) error
}

type User struct {
	ID             int
	Username       string
	Email          string
	HashedPassword []byte
	Administrator  bool
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(username, email, password string, administrator bool) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, administrator)
    VALUES(?, ?, ?, ?)`

	_, err = m.DB.Exec(stmt, username, email, string(hashedPassword), administrator)
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

func (m *UserModel) Authenticate(emailOrUsername, password string) (int, error) {
	var id int
	var hashedPassword []byte

	// Try to find user by email first, then by username (name field in database)
	stmt := "SELECT id, hashed_password FROM users WHERE email = ? OR name = ?"

	err := m.DB.QueryRow(stmt, emailOrUsername, emailOrUsername).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"

	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}

func (m *UserModel) IsAdministrator(id int) (bool, error) {
	var isAdmin bool

	stmt := "SELECT administrator FROM users WHERE id = ?"

	err := m.DB.QueryRow(stmt, id).Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return isAdmin, nil
}

func (m *UserModel) ChangePassword(id int, currentPassword, newPassword string) error {
	var hashedPassword []byte

	// First, get the current password hash for verification
	stmt := "SELECT hashed_password FROM users WHERE id = ?"
	err := m.DB.QueryRow(stmt, id).Scan(&hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidCredentials
		}
		return err
	}

	// Verify the current password
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(currentPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		}
		return err
	}

	// Hash the new password
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}

	// Update the password in the database
	stmt = "UPDATE users SET hashed_password = ? WHERE id = ?"
	_, err = m.DB.Exec(stmt, string(newHashedPassword), id)
	return err
}
