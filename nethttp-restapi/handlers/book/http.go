package book

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"zopsmart.com/nethttp-test/models"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case "GET":
		a_id := r.PathValue("a_id")
		b_id := r.PathValue("b_id")

		parsedAuthorID, err := strconv.ParseInt(a_id, 10, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "failed to parse author id")
			return
		}

		parsedBookID, err := strconv.ParseInt(b_id, 10, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "failed to parse book id")
			return
		}

		data, err := models.GetBook(parsedAuthorID, parsedBookID)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "book not found")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)

		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "Failed to encode data as JSON", http.StatusInternalServerError)
		}
	case "PUT":
		var a models.Author

		err := json.NewDecoder(r.Body).Decode(&a)
		if err != nil {
			fmt.Println(err)
		}

		io.WriteString(w, "Book Updated")
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func HandlerCreate(w http.ResponseWriter, r *http.Request) {

	switch method := r.Method; method {
	case "POST":
		a_id := r.PathValue("a_id")

		parsedAuthorID, err := strconv.ParseInt(a_id, 10, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "failed to parse author id")
			return
		}

		data, err := models.GetAuthor(parsedAuthorID)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "author not found")
			return
		}

		var b models.Book

		err = json.NewDecoder(r.Body).Decode(&b)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "failed to parse the data")
			return
		}

		b.AuthorID = data.ID
		err = b.Save()

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "error saving to database")
			return
		}

		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, "Book created")
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
