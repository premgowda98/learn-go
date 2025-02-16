package models

import (
	"errors"

	"github.com/google/uuid"
)

type Engine struct {
	ID           uuid.UUID `json:"id"`
	Displacement int64     `json:"name"`
	Cyclinders   int64     `json:"cylinders"`
	Range        int64     `json:"range"`
}

type EngineRequest struct {
	Displacement int64     `json:"name"`
	Cyclinders   int64     `json:"cylinders"`
	Range        int64     `json:"range"`
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

func validateDisplacement(displacement int64) error {
	if displacement <= 0 {
		return errors.New("displacement must be greater then 0")
	}

	return nil
}

func validateCylinders(cylinders int64) error {
	if cylinders <= 0 {
		return errors.New("cylinders must be greater then 0")
	}

	return nil
}

func validateRange(rangeK int64) error {
	if rangeK <= 0 {
		return errors.New("range must be greater then 0")
	}

	return nil
}
