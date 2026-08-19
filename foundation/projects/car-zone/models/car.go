package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Car struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Year      string    `json:"year"`
	Brand     string    `json:"brand"`
	FuelType  string    `json:"fuel_type"`
	Engine    Engine    `json:"engine"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CarRequest struct {
	Name     string  `json:"name"`
	Year     string  `json:"year"`
	Brand    string  `json:"brand"`
	FuelType string  `json:"fuel_type"`
	Engine   Engine  `json:"engine"`
	Price    float64 `json:"price"`
}

func (c *CarRequest) Validate() error {
	if err := validateName(c.Name); err != nil {
		return err
	}
	if err := validateBrand(c.Brand); err != nil {
		return err
	}
	if err := validateYear(c.Year); err != nil {
		return err
	}
	if err := validateFuelType(c.FuelType); err != nil {
		return err
	}
	if err := validatePrice(c.Price); err != nil {
		return err
	}

	return nil
}

func validateName(name string) error {
	if name == "" {
		return errors.New("name is requried")
	}

	return nil
}

func validateYear(year string) error {
	if year == "" {
		return errors.New("year is requried")
	}

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return errors.New("year should be valid number")
	}

	currentYear := time.Now().Year()

	if yearInt < 1988 || yearInt > currentYear {
		return errors.New("year must be between 1886 and current year")
	}

	return nil
}

func validateBrand(brand string) error {
	if brand == "" {
		return errors.New("brand is requried")
	}

	return nil
}

func validateFuelType(fuelType string) error {
	validFuel := []string{"petrol", "diesel", "cng", "electric"}

	if fuelType == "" {
		return errors.New("fuel_type is requried")
	}

	isValid := false
	for _, valid := range validFuel {
		if fuelType == valid {
			isValid = true
			break
		}
	}

	if !isValid {
		return errors.New("fuel_type is invalid")
	}
	return nil
}

func validatePrice(price float64) error {
	if price <= 0 {
		return errors.New("price should be greater than 0")
	}

	return nil
}
