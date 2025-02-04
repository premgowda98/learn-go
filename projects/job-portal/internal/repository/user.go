package repository

import (
	"database/sql"
	"project/user-management/internal/models"
)

func CreateUser(db *sql.DB, user *models.User) error {
	_, err := db.Exec(`INSERT INTO users (username, password, email, is_admin) VALUES (?,?,?,?)`, user.Username, user.Password, user.Email, user.IsAdmin)
	return err
}

func GetUserByID(db *sql.DB, id int) (*models.User, error) {
	row := db.QueryRow(`SELECT id, username, email, is_admin, created_at, updated_at FROM users WHERE id=?`, id)

	user := &models.User{}

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByUsername(db *sql.DB, username *string) (*models.User, error) {
	row := db.QueryRow(`SELECT id, username, password FROM users WHERE username=?`, username)
	user := &models.User{}

	err := row.Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		return nil, err
	}

	return user, nil
}
