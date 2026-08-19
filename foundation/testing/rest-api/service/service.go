package service

import (
	"errors"
	"fmt"
	"test/restapi/models"
)

type Service struct {
	store UserStore
}

func New(store UserStore) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) Create(u *models.UserRequest) error {
	fmt.Println("Recieved from handler", u.Name)

	if u.Name == "fail" {
		return errors.New("not allowed")
	}

	model_user := models.User{
		ID:   u.ID,
		Name: u.Name,
	}

	err := s.store.Create(&model_user)

	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Get(id int) (*models.User, error){
	user, err := s.store.Get(id)

	if err !=nil {
		return nil, err
	}

	return user, nil
}
