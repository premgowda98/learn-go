package store

import (
	"test/restapi/models"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestShouldCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal("error opening stub data connection")
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO users").WithArgs("ko").WillReturnResult(sqlmock.NewResult(1,1))
	mock.ExpectCommit()
	store := Store{db: db}

	user := models.User{
		ID: 0,
		Name: "ko",
	}
	if err = store.Create(&user); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}
