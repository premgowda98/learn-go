package car

import (
	"encoding/json"
	"io"
	"net/http"
	"project/car-zone/models"
	"project/car-zone/service"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	service service.Car
}

func New(service service.Car) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetCarById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := h.service.GetCarById(ctx, id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(&res)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	_, _ = w.Write(body)

}

func (h *Handler) GetCarByBrand(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	brand := r.URL.Query().Get("brand")

	res, err := h.service.GetCarByBrand(ctx, brand)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(&res)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	_, _ = w.Write(body)

}

func (h *Handler) CreateCar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var carRequest *models.CarRequest

	err = json.Unmarshal(body, &carRequest)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	car, err := h.service.CreateCar(ctx, carRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respBody, err := json.Marshal(&car)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_, _ = w.Write(respBody)
}
