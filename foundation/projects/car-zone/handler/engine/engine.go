package engine

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
	service service.Engine
}

func New(service service.Engine) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetEngineById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := h.service.GetEngineById(ctx, id)

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

func (h *Handler) CreateEngine(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var engineRequest *models.EngineRequest

	err = json.Unmarshal(body, &engineRequest)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	car, err := h.service.CreateEngine(ctx, engineRequest)
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
