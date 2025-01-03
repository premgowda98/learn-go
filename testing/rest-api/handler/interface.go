package handler

import "test/restapi/models"

type UserService interface{
	Create(u *models.UserRequest) error
	Get(id int)(*models.User, error)
}