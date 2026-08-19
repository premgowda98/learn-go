package store

import (
	"database/sql"
	"fmt"
	"test/restapi/models"
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Create(u *models.User) error {
	fmt.Println("Created user", u.Name)

	query := `INSERT INTO users (name) VALUES (?)`

	sql_smt, err := s.db.Prepare(query)

	if err != nil {
		return err
	}

	defer sql_smt.Close()

	_, err = sql_smt.Exec(u.Name)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Get(id int) (*models.User, error) {
	query := `SELECT id, name FROM users WHERE id=?`
	rows := s.db.QueryRow(query, id)

	var user models.User

	err := rows.Scan(&user.ID, &user.Name)

	if err != nil {
		return nil, err
	}

	return &user, nil

}
