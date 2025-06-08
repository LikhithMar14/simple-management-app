package service

import (
	"context"

	"github.com/LikhithMar14/management/models"
)

type CarService interface {
	GetCarByID(ctx context.Context, id string) (models.Car, error)
	GetCarsByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error)
	CreateCar(ctx context.Context, car *models.CarRequest) (models.Car, error)
	UpdateCar(ctx context.Context, id string, car *models.CarRequest) (models.Car, error)
	DeleteCar(ctx context.Context, id string) error
}

type EngineService interface {
	GetEngineByID(ctx context.Context, id string) (models.Engine, error)
	CreateEngine(ctx context.Context, engine *models.EngineRequest) (models.Engine, error)
	UpdateEngine(ctx context.Context, id string, engine *models.EngineRequest) (models.Engine, error)
	DeleteEngine(ctx context.Context, id string) error
}