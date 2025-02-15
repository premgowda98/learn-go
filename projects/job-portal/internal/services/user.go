package services

import (
	"database/sql"
	"project/user-management/internal/models"
	"project/user-management/internal/repository"
)

func GetUserByID(db *sql.DB, id int) (*models.User, error) {
	user, err := repository.GetUserByID(db, id)

	if err != nil {
		return nil, err
	}

	return user, nil

}

func UpdateUser(db *sql.DB, user *models.User, id int) (*models.User, error) {
	user, err := repository.UpdateUser(db, user, id)

	if err != nil {
		return nil, err
	}

	return user, nil
}
