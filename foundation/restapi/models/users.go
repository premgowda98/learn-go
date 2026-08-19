package models

import (
	"errors"
	"project/restapi/db"
	"project/restapi/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u User) Save() error {
	query := `INSERT INTO users (email, password) VALUES (?,?)`
	sql_smt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer sql_smt.Close()

	hashedPassword, err := utils.HashPassword(u.Password)

	if err != nil {
		return err
	}

	_, err = sql_smt.Exec(u.Email, hashedPassword)

	if err != nil {
		return err
	}

	return err
}

func (u *User) ValidateCredentials() error {
	query := `SELECT id, password FROM users WHERE email = ?`

	user := db.DB.QueryRow(query, u.Email)
	
	var retrievedPassword string
	err := user.Scan(&u.ID, &retrievedPassword)

	if err != nil{
		return nil
	}

	passwordIsValid := utils.CheckHashPassword(retrievedPassword, u.Password)

	if !passwordIsValid{
		return errors.New("invalid credentials")
	}

	return nil
}
