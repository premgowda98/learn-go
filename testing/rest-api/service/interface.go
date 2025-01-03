package service

import "test/restapi/models"

type UserStore interface{
	Create(u *models.User) error
	Get(id int) (*models.User, error)
}