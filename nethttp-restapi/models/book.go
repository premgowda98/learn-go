package models

import (
	"time"

	"zopsmart.com/nethttp-test/db"
)

type Book struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	ISBN        string    `json:"isbn"`
	AuthorID    int64     `json:"author_id"`
	ReleaseDate time.Time `json:"release_date"`
}

func (b *Book) New(name string, isbn string, release_date time.Time, author_id int64) (Book, error) {
	return Book{
		Name:        name,
		ISBN:        isbn,
		ReleaseDate: release_date,
		AuthorID:    author_id,
	}, nil
}

func (b *Book) Save() error {
	query := `INSERT INTO books (name, isbn, release_date, author_id)
	VALUES (?,?,?,?)
	`
	sql_smt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer sql_smt.Close()

	result, err := sql_smt.Exec(b.Name, b.ISBN, b.ReleaseDate, b.AuthorID)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	b.ID = id

	return err
}

func GetBook(a_id int64, b_id int64) (*Book, error) {
	query := `SELECT id, name, isbn, release_date, author_id FROM books WHERE id=? AND author_id=?`

	rows := db.DB.QueryRow(query, b_id, a_id)

	var book Book

	err := rows.Scan(&book.ID, &book.Name, &book.ISBN, &book.ReleaseDate, &book.AuthorID)

	if err != nil {
		return nil, err
	}

	return &book, nil
}
