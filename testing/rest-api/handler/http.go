package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"test/restapi/models"
)

type Handler struct {
	service UserService
}

func New(service UserService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {

	var user models.UserRequest

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = h.service.Create(&user)

	if err != nil {
		fmt.Println(err)
	}

	io.WriteString(w, "created")
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("recieved get user")
	id := r.PathValue("id")

	parsedAuthorID, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		fmt.Println("dailed to parse path value")
	}

	rows, err := h.service.Get(int(parsedAuthorID))

	if err != nil {
		fmt.Println("could not fetch values")
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(rows); err != nil {
		http.Error(w, "Failed to encode data as JSON", http.StatusInternalServerError)
	}
}
