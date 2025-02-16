package engine

import (
	"context"
	"project/car-zone/models"
	"project/car-zone/store"
)

type Service struct {
	store store.Engine
}

func New(store store.Engine) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) GetEngineById(ctx context.Context, id int64) (*models.Engine, error) {
	engine, err := s.store.GetEngineById(ctx, id)

	if err != nil {
		return nil, err
	}

	return engine, nil
}

func (s *Service) CreateEngine(ctx context.Context, engineReq *models.EngineRequest) (*models.Engine, error) {
	if err := engineReq.Validate(); err != nil {
		return nil, err
	} 

	engine, err := s.store.CreateEngine(ctx, engineReq)

	if err != nil {
		return nil, err
	}

	return engine, nil
}
