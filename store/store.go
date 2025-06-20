package store

import (
	"context"
	"database/sql"

	"github.com/LikhithMar14/management/models"
	"github.com/LikhithMar14/management/store/car"
	"github.com/LikhithMar14/management/store/engine"
)

type Storage struct {
	CarStore CarStoreInterface
	EngineStore EngineStoreInterface
}

type CarStoreInterface interface {
	GetCarByID(ctx context.Context, id string) (models.Car, error)
	GetCarsByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error)
	CreateCar(ctx context.Context, car *models.CarRequest) (models.Car, error)
	UpdateCar(ctx context.Context, id string, car *models.CarRequest) (models.Car, error)
	DeleteCar(ctx context.Context, id string) error
}

type EngineStoreInterface interface {
	GetEngineByID(ctx context.Context, id string) (models.Engine, error)
	CreateEngine(ctx context.Context, engine *models.EngineRequest) (models.Engine, error)
	UpdateEngine(ctx context.Context, id string, engine *models.EngineRequest) (models.Engine, error)
	DeleteEngine(ctx context.Context, id string) error
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{	
		CarStore: car.NewCarStore(db),
		EngineStore: engine.NewEngineStore(db),
	}
}