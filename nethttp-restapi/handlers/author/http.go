package author

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
		id := r.PathValue("id")

		parsedID, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "failed to parse route")
			return
		}

		data, err := models.GetAuthor(parsedID)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "author not found")
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

		// err = a.Save()

		// if err != nil {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	io.WriteString(w, "error saving to database")
		// }

		// w.WriteHeader(http.StatusCreated)
		io.WriteString(w, "Author Updated")
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func HandlerCreate(w http.ResponseWriter, r *http.Request) {

	switch method := r.Method; method {
	case "POST":
		var a models.Author

		err := json.NewDecoder(r.Body).Decode(&a)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "failed to parse the data")
			return
		}

		err = a.Save()

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "error saving to database")
			return
		}

		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, "Author created")
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
