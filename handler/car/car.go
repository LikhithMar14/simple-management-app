package car

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/LikhithMar14/management/models"
	"github.com/LikhithMar14/management/service"
	"github.com/go-chi/chi/v5"
)

type CarHandler struct {
	service service.CarService
}

func NewCarHandler(service service.CarService) *CarHandler {
	return &CarHandler{service: service}
}

func (h *CarHandler) GetCarByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := chi.URLParam(r, "id")
	log.Print("I AM IN CAR BY ID HANDLER")

	car, err := h.service.GetCarByID(ctx, vars)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(car)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (h *CarHandler) GetCarsByBrand(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	brand := r.URL.Query().Get("brand")
	isEngine := r.URL.Query().Get("isEngine") == "true"

	log.Println("I am Get cars by brand:", brand)

	cars, err := h.service.GetCarsByBrand(ctx, brand, isEngine)
	if err != nil {
		http.Error(w, "Failed to fetch cars", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cars); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}


func (h *CarHandler) CreateCar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var newCar models.CarRequest
	err := json.NewDecoder(r.Body).Decode(&newCar)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	car , err := h.service.CreateCar(ctx, &newCar)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(car)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}


func (h *CarHandler) UpdateCar(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()

	var updateCar models.CarRequest

	err := json.NewDecoder(r.Body).Decode(&updateCar)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vars := chi.URLParam(r, "id")
	fmt.Print("id: ",vars)

	car, err := h.service.UpdateCar(ctx, vars, &updateCar)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("error: ",err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(car)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (h *CarHandler) DeleteCar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := chi.URLParam(r, "id")

	err := h.service.DeleteCar(ctx, vars)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}


