package mocks

import (
	"scvis.vaccsedu.org/internal/models"
)

type UserModel struct{}

func (m *UserModel) Insert(username, email, password string, administrator bool) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(emailOrUsername, password string) (int, error) {
	if (emailOrUsername == "alice@example.com" || emailOrUsername == "alice") && password == "pa$$word" {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}

func (m *UserModel) IsAdministrator(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}

func (m *UserModel) ChangePassword(id int, currentPassword, newPassword string) error {
	if id == 1 && currentPassword == "pa$$word" {
		return nil
	}
	return models.ErrInvalidCredentials
}
