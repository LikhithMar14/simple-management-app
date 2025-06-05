package models

import (
	"errors"

	"github.com/google/uuid"
)

type Engine struct {
	EngineID uuid.UUID `json:"engine_id"`
	Displacement int64 `json:"displacement"`
	NumberOfCylinders int64 `json:"number_of_cylinders"`
	CarRange int64 `json:"car_range"`
}

type EngineRequest struct {
	Displacement int64 `json:"displacement"`
	NumberOfCylinders int64 `json:"number_of_cylinders"`
	CarRange int64 `json:"car_range"`
}	

func ValidateEngineRequest(engine EngineRequest) error {
	if err := validateDisplacement(engine.Displacement); err != nil {
		return err
	}
	if err := validateNumberOfCylinders(engine.NumberOfCylinders); err != nil {
		return err
	}
	if err := validateCarRange(engine.CarRange); err != nil {
		return err
	}
	return nil
}

func validateDisplacement(displacement int64) error {
	if displacement <= 0 {
		return errors.New("displacement must be greater than 0")
	}
	return nil
}

func validateNumberOfCylinders(numberOfCylinders int64) error {
	if numberOfCylinders <= 0 {
		return errors.New("number of cylinders must be greater than 0")
	}
	return nil	
}

func validateCarRange(carRange int64) error {
	if carRange <= 0 {
		return errors.New("car range must be greater than 0")
	}
	return nil
}