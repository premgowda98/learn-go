package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"test/restapi/handler"
	"test/restapi/service"
	"test/restapi/store"

	_ "github.com/mattn/go-sqlite3"
)

func main(){
	fmt.Println("Starting server")

	db, err := sql.Open("sqlite3", "sample.db")

	if err !=nil {
		fmt.Println("Database connection failed")
		return
	}

	createTable := `CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT)`

	_, err = db.Exec(createTable)

	if err !=nil {
		fmt.Println("failed to create database")
	}

	str := store.New(db)
	service := service.New(str)
	handler := handler.New(service)

	http.HandleFunc("POST /user", handler.CreateUser)
	http.HandleFunc("GET /user/{id}", handler.GetUser)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "healthy")
	})
	fmt.Println("Server running at port 8090")
	http.ListenAndServe(":8009", nil)
}