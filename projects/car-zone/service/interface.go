package service

import (
	"context"
	"project/car-zone/models"
)

type Car interface {
	GetCarById(ctx context.Context, id int64) (*models.Car, error)
	GetCarByBrand(ctx context.Context, brand string) ([]*models.Car, error)
	CreateCar(ctx context.Context, carReq *models.CarRequest) (*models.Car, error)
}

type Engine interface {
	GetEngineById(ctx context.Context, id int64) (*models.Engine, error)
	CreateEngine(ctx context.Context, engineReq *models.EngineRequest) (*models.Engine, error)
}
