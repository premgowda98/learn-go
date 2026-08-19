package store

import (
	"context"
	"project/car-zone/models"
)

type Car interface {
	GetCarById(ctx context.Context, id int) (*models.Car, error)
	GetCarByBrand(ctx context.Context, brand string) ([]*models.Car, error)
	CreateCar(ctx context.Context, carRequest *models.CarRequest) (*models.Car, error)
	DeleteCar(ctx context.Context, id int) error
}

type Engine interface {
	GetEngineById(ctx context.Context, id int) (*models.Engine, error)
	CreateEngine(ctx context.Context, engineRequest *models.EngineRequest) (*models.Engine, error)
	DeleteEngine(ctx context.Context, id int) error
}
