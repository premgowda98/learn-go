package models

import "zopsmart.com/nethttp-test/db"

type Author struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Age     int32  `json:"age"`
	Address string `json:"address"`
}

func (a *Author) New(name string, age int32, address string) (Author, error) {
	return Author{
		Name:    name,
		Age:     age,
		Address: address,
	}, nil
}

func (a *Author) Save() error {
	query := `INSERT INTO authors (name, age, address)
	VALUES (?,?,?)
	`
	sql_smt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer sql_smt.Close()

	result, err := sql_smt.Exec(a.Name, a.Age, a.Address)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	a.ID = id

	return err
}

func GetAuthor(id int64) (*Author, error) {
	query := `SELECT id, name, age, address FROM authors where id=?`

	rows := db.DB.QueryRow(query, id)

	var author Author

	err := rows.Scan(&author.ID, &author.Name, &author.Age, &author.Address)

	if err != nil {
		return nil, err
	}

	return &author, nil
}
