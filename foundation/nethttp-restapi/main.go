package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/joho/godotenv"

	"zopsmart.com/nethttp-test/db"
	"zopsmart.com/nethttp-test/handlers/author"
	"zopsmart.com/nethttp-test/handlers/book"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("failed to laod env file")
	}
	db.InitDB()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("server healthy")
		io.WriteString(w, "server healthy")
	})

	http.HandleFunc("/author", author.HandlerCreate)
	http.HandleFunc("/author/{id}", author.Handler)
	http.HandleFunc("/author/{a_id}/book", book.HandlerCreate)
	http.HandleFunc("/author/{a_id}/book/{b_id}", book.Handler)

	fmt.Println("Running http server on port 8000")
	err = http.ListenAndServe(":8000", nil)

	if err != nil {
		fmt.Println("HTTP server closed")
	}

}
