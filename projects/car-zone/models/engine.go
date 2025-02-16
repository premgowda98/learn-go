package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Engine struct {
	ID           uuid.UUID `json:"id"`
	Displacement int       `json:"name"`
	Cyclinders   int       `json:"cylinders"`
	Range        int       `json:"range"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type EngineRequest struct {
	Displacement int `json:"name"`
	Cyclinders   int `json:"cylinders"`
	Range        int `json:"range"`
}

func (e *EngineRequest) Validate() error {
	if err := validateDisplacement(e.Displacement); err != nil {
		return err
	}

	if err := validateCylinders(e.Cyclinders); err != nil {
		return err
	}

	if err := validateRange(e.Range); err != nil {
		return err
	}

	return nil
}

func validateDisplacement(displacement int) error {
	if displacement <= 0 {
		return errors.New("displacement must be greater then 0")
	}

	return nil
}

func validateCylinders(cylinders int) error {
	if cylinders <= 0 {
		return errors.New("cylinders must be greater then 0")
	}

	return nil
}

func validateRange(rangeK int) error {
	if rangeK <= 0 {
		return errors.New("range must be greater then 0")
	}

	return nil
}
