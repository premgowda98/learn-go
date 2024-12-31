package main

import (
	"log"
	"os"
	"project/restapi-2/config"
	"project/restapi-2/db"

	mysqlConfig "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
)

func main() {
	db, err := db.NewMySQLStorage(mysqlConfig.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		DBName:               config.Envs.DBName,
		Addr:                 config.Envs.DBAddress,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})

	if err != nil {
		log.Fatal(err)
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})

	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"mysql",
		driver,
	)

	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[(len(os.Args) - 1)]

	if cmd == "up" {
		if err = m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
	if cmd == "down" {
		if err = m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}

}
