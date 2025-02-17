package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	conString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", conString)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	err = createTable(db)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func createTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS engine (
		id UUID PRIMARY KEY,
		displacement INT NOT NULL,
		cylinders INT NOT NULL,
		range INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS car (
		id UUID PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		year VARCHAR(4) NOT NULL,
		brand VARCHAR(255) NOT NULL,
		fuel_type VARCHAR(50) NOT NULL,
		engine UUID NOT NULL,
		price DECIMAL(10, 2) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)

	if err != nil {
		return err
	}

	return nil
}
