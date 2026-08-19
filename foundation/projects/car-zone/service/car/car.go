package car

import (
	"context"
	"project/car-zone/models"
	"project/car-zone/store"
)

type Service struct {
	store store.Car
}

func New(store store.Car) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) GetCarById(ctx context.Context, id int) (*models.Car, error) {
	car, err := s.store.GetCarById(ctx, id)

	if err != nil {
		return nil, err
	}

	return car, nil
}

func (s *Service) GetCarByBrand(ctx context.Context, brand string) ([]*models.Car, error) {
	car, err := s.store.GetCarByBrand(ctx, brand)

	if err != nil {
		return nil, err
	}

	return car, nil
}

func (s *Service) CreateCar(ctx context.Context, carReq *models.CarRequest) (*models.Car, error) {

	if err := carReq.Validate(); err != nil {
		return nil, err
	}

	car, err := s.store.CreateCar(ctx, carReq)

	if err != nil {
		return nil, err
	}

	return car, nil
}
