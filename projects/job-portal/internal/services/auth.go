package services

import (
	"database/sql"
	"project/user-management/internal/models"
	"project/user-management/internal/repository"
	"project/user-management/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

func LoginUser(db *sql.DB, username string, password string) (string, error) {
	user, err := repository.GetUserByUsername(db, &username)

	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password))

	if err != nil {
		return "", err
	}

	return utils.GenerateJWTToken(username)
}

func RegisterUser(db *sql.DB, user *models.User) error {
	hashedPassword, err := HashPassword(user.Password)

	if err != nil {
		return err
	}

	user.Password = &hashedPassword

	return repository.CreateUser(db, user)
}

func HashPassword(password *string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
