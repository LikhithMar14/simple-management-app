package engine

import (
	"encoding/json"
	"net/http"

	"github.com/LikhithMar14/management/models"
	"github.com/LikhithMar14/management/service"
	"github.com/go-chi/chi/v5"
)


type EngineHandler struct {
	service service.EngineService
}

func NewEngineHandler(service service.EngineService) *EngineHandler {
	return &EngineHandler{service: service}
}

func (h *EngineHandler) GetEngineByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := chi.URLParam(r, "id")

	engine, err := h.service.GetEngineByID(ctx, vars)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(engine)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func (h *EngineHandler) CreateEngine(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var newEngine models.EngineRequest
	err := json.NewDecoder(r.Body).Decode(&newEngine)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	engine, err := h.service.CreateEngine(ctx, &newEngine)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(engine)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *EngineHandler) UpdateEngine(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := chi.URLParam(r, "id")

	var updateEngine models.EngineRequest
	err := json.NewDecoder(r.Body).Decode(&updateEngine)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	engine, err := h.service.UpdateEngine(ctx, vars, &updateEngine)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(engine)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *EngineHandler) DeleteEngine(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := chi.URLParam(r, "id")

	err := h.service.DeleteEngine(ctx, vars)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]string{"message": "Engine deleted successfully"})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

