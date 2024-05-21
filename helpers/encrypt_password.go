package helpers

import (
	"golang.org/x/crypto/bcrypt"
	"jwt-authentication/models"
)

func Encrypt(user *models.User) ([]byte, error) {

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}
	return hashPassword, nil
}
