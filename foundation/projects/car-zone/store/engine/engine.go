package engine

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

func (s *Store) GetEngineById(ctx context.Context, id int) (*models.Engine, error) {
	var engine models.Engine

	query := `SELECT e.id, e.displacement, e.cyclinders, e.range FROM engine e  WHERE id=$1`

	rows := s.db.QueryRowContext(ctx, query, id)

	err := rows.Scan(&engine.ID, &engine.Displacement, &engine.Cyclinders, &engine.Range)

	if err != nil {
		if err == sql.ErrNoRows {
			return &engine, nil
		}
		return nil, err
	}

	return &engine, nil
}

func (s *Store) CreateEngine(ctx context.Context, engineRequest *models.EngineRequest) (*models.Engine, error) {
	createdEngine := &models.Engine{}
	engineId := uuid.New()
	createdAt := time.Now()

	newEngine := &models.Engine{
		ID:           engineId,
		Displacement: engineRequest.Displacement,
		Cyclinders:   engineRequest.Cyclinders,
		Range:        engineRequest.Range,
		CreatedAt:    createdAt,
	}

	query := `INSERT INTO engine (id, displacement, cyclinders, ranges, created_at)
	VALUES ($1,$2,$3,$4,$5) RETURNING id, displacement, cyclinders, ranges, created_at`

	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return createdEngine, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	err = tx.QueryRowContext(ctx, query, &newEngine.ID, &newEngine.Displacement, &newEngine.Cyclinders, &newEngine.Range,
		&newEngine.CreatedAt).Scan(&createdEngine.ID, &createdEngine.Displacement, &createdEngine.Cyclinders,
		&createdEngine.Range, &createdEngine.CreatedAt)

	if err != nil {
		return createdEngine, err
	}

	return createdEngine, nil

}

func (s *Store) DeleteEngine(ctx context.Context, id int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	result, err := tx.ExecContext(ctx, `DELETE FROM engine WHERE id=$1`, id)

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
