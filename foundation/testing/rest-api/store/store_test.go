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

	mock.ExpectPrepare(`INSERT INTO users \(name\) VALUES \(\?\)`).WillReturnError(nil)
	mock.ExpectExec(`INSERT INTO users \(name\) VALUES \(\?\)`).WithArgs("ko").WillReturnResult(sqlmock.NewResult(1, 1))

	store := Store{db: db}

	user := models.User{
		ID:   0,
		Name: "ko",
	}
	if err = store.Create(&user); err != nil {
		t.Errorf("error was not expected while inserting records: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestShouldGetUser(t *testing.T){
	db, mock, err := sqlmock.New()

	if err !=nil {
		t.Fatal("error opening stub connection")
	}

	defer db.Close()

	store := Store{db:db}

	expectedUser := &models.User{
        ID:   1,
        Name: "John Doe",
    }

	mock.ExpectQuery(`SELECT id, name FROM users WHERE id=\?`).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(expectedUser.ID, expectedUser.Name))

	if _, err = store.Get(1); err !=nil{
		t.Fatal("not expecting error")
	}
}
