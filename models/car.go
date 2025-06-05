package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
	"slices"
)

type Car struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Year      string    `json:"year"`
	Brand     string    `json:"brand"`
	FuelType  string    `json:"fuel_type"`
	Engine    Engine    `json:"engine"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CarRequest struct {
	Name      string    `json:"name"`
	Year      string    `json:"year"`
	Brand     string    `json:"brand"`
	FuelType  string    `json:"fuel_type"`
	Engine    Engine    `json:"engine"`
	Price     float64   `json:"price"`
}

func ValidateCarRequest(car CarRequest) error {
	if err := validateName(car.Name); err != nil {
		return err
	}
	if err := validateYear(car.Year); err != nil {
		return err
	}
	if err := validateBrand(car.Brand); err != nil {
		return err
	}
	if err := validateFuelType(car.FuelType); err != nil {
		return err
	}
	if err := validateEngine(car.Engine); err != nil {
		return err	
	}
	if err := validatePrice(car.Price); err != nil {
		return err
	}
	return nil
}


func validateName(name string) error {
	if name == "" {
		return errors.New("name is required")
	}
	return nil
}

func validateYear(year string) error {
	if year == "" {
		return errors.New("year is required")
	}
	if len(year) != 4 {
		return errors.New("year must be 4 digits")
	}
	value,err := strconv.Atoi(year)
	if err != nil {
		return errors.New("year must be a number")
	}
	if value < 1900 || value > time.Now().Year() {
		return errors.New("year must be between 1900 and current year")
	}
	return nil
}

func validateBrand(brand string) error {
	if brand == "" {
		return errors.New("brand is required")
	}
	return nil
}

func validateFuelType(fuelType string) error {
	validFuelTypes := []string{"petrol", "diesel", "electric", "hybrid"}
	if slices.Contains(validFuelTypes, fuelType) {
			return nil
		}
	return errors.New("fuel type must be one of the following: petrol, diesel, electric, hybrid")
}

func validateEngine(engine Engine) error {
	if engine.Displacement <= 0 {
		return errors.New("displacement must be greater than 0")
	}																																																	
	if engine.NumberOfCylinders <= 0 {
		return errors.New("number of cylinders must be greater than 0")
	}
	if engine.EngineID == uuid.Nil {
		return errors.New("engine id is required")
	}
	if engine.Displacement < 1000 {
		return errors.New("displacement must be greater than 1000")
	}
	if engine.CarRange <= 0 {
		return errors.New("car range must be greater than 0")
	}
	return nil
}

func validatePrice(price float64) error {
	if price <= 0 {
		return errors.New("price must be greater than 0")
	}
	return nil
}
