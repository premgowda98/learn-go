package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	fmt.Println("Initailizing database")

	var err error
	DB, err = sqliteDB()
	// DB, err = mysqlDB()

	if err != nil {
		panic("could not connect to database")
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxOpenConns(5)

	createTables()
}

func sqliteDB() (*sql.DB, error) {
	var err error
	DB, err = sql.Open("sqlite3", "api.db")

	if err != nil {
		panic("could not connect to database")
	}
	return DB, nil
}

func mysqlDB() (*sql.DB, error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	conn_str := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, host, dbName)

	DB, err := sql.Open("mysql", conn_str)

	if err != nil {
		panic("could not connect to database")
	}
	return DB, nil
}

func createTables() {
	booksTable := `
	CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		isbn TEXT NOT NULL,
		release_date DATETIME,
		author_id INTEGER
	)
	`

	// booksTable := `
	// CREATE TABLE IF NOT EXISTS bookmanagement.books (
	// 	id INT AUTO_INCREMENT PRIMARY KEY,
	// 	name VARCHAR(255) NOT NULL,
	// 	isbn VARCHAR(20) NOT NULL,
	// 	release_date DATETIME,
	// 	author_id INT
	// );
	// `

	_, err := DB.Exec(booksTable)

	if err != nil {
		panic("could not create books table")
	}

	authorsTable := `
	CREATE TABLE IF NOT EXISTS authors (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		age INTEGER,
		address TEXT
	)
	`

	// authorsTable := `
	// CREATE TABLE IF NOT EXISTS bookmanagement.authors (
	// 	id INT AUTO_INCREMENT PRIMARY KEY ,
	// 	name VARCHAR(250),
	// 	age INT,
	// 	address VARCHAR(250)
	// )
	// `
	_, err = DB.Exec(authorsTable)

	if err != nil {
		panic("could not create authors table")
	}
}
