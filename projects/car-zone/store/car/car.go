package car

import (
	"context"
	"database/sql"
	"errors"
	"project/car-zone/models"
	"time"

	"github.com/google/uuid"
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetCarById(ctx context.Context, id int) (*models.Car, error) {
	var car models.Car

	query := `SELECT c.id, c.name, c.year, c.brand, c.fuel_type, c.price, c.created_at, c.updated_at,
	e.id, e.displacement, e.cyclinders, e.range FROM cars c LEFT JOIN engine e on c.engine_id=c.id WHERE id=$1`

	rows := s.db.QueryRowContext(ctx, query, id)

	err := rows.Scan(&car.ID, &car.Name, &car.Year, &car.Brand, &car.FuelType, &car.Price, &car.CreatedAt, &car.UpdatedAt,
		&car.Engine.ID, &car.Engine.Displacement, &car.Engine.Cyclinders, &car.Engine.Range)

	if err != nil {
		if err == sql.ErrNoRows {
			return &car, nil
		}
		return nil, err
	}

	return &car, nil
}

func (s *Store) GetCarByBrand(ctx context.Context, brand string) ([]*models.Car, error) {
	var cars []*models.Car

	query := `SELECT c.id, c.name, c.year, c.brand, c.fuel_type, c.price, c.created_at, c.updated_at FROM cars c WHERE brand=$1`

	rows, err := s.db.QueryContext(ctx, query, brand)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var car models.Car
		err := rows.Scan(&car.ID, &car.Name, &car.Year, &car.Brand, &car.FuelType, &car.Price, &car.CreatedAt, &car.UpdatedAt)
		if err != nil {
			return nil, err
		}

		cars = append(cars, &car)
	}

	return cars, nil
}

func (s *Store) CreateCar(ctx context.Context, carRequest *models.CarRequest) (*models.Car, error) {
	createdCar := &models.Car{}
	var engineId uuid.UUID

	err := s.db.QueryRowContext(ctx, `SELECT id from engine where id=$1`, &carRequest.Engine.ID).Scan(&engineId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return createdCar, errors.New("engine id not present")
		}
		return createdCar, err
	}

	carId := uuid.New()
	createdAt := time.Now()

	newCar := &models.Car{
		ID:        carId,
		Name:      carRequest.Name,
		Brand:     carRequest.Brand,
		Year:      carRequest.Year,
		FuelType:  carRequest.FuelType,
		Price:     carRequest.Price,
		Engine:    carRequest.Engine,
		CreatedAt: createdAt,
	}

	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return createdCar, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	query := `INSERT INTO car (id, name, year, brand, fuel_type, engine_id, price, created_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8)
	RETURNING id, name, year, brand, fuel_type, engine_id, price, created_at`

	err = tx.QueryRowContext(ctx, query, &newCar.ID, &newCar.Name, &newCar.Year, &newCar.Brand,
		&newCar.FuelType, &newCar.Engine.ID, &newCar.Price, &newCar.CreatedAt).Scan(
		&createdCar.ID, &createdCar.Name, &createdCar.Year, &createdCar.Brand, &createdCar.FuelType,
		&createdCar.Engine.ID, &createdCar.Price, &createdCar.CreatedAt,
	)

	if err != nil {
		return createdCar, err
	}

	return createdCar, nil

}

func (s *Store) DeleteCar(ctx context.Context, id int) error {

	tx, err := s.db.BeginTx(ctx, nil)
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	result, err := tx.ExecContext(ctx, `DELETE FROM car WHERE id=$1`, id)

	if err != nil {
		return err
	}

	rowsEffect, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsEffect == 0 {
		return errors.New("no records found")
	}

	return nil

}
