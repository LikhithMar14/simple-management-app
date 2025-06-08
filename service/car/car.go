package car

import (
	"context"
	"fmt"
	"log"

	"github.com/LikhithMar14/management/models"
	"github.com/LikhithMar14/management/store"
)


type CarService struct {
	store store.CarStoreInterface
}

func NewCarService(store store.CarStoreInterface) *CarService {
	return &CarService {
		store: store,
	}
}

func (s *CarService) GetCarByID(ctx context.Context, id string) (models.Car, error) {
	log.Println("I am in car service",id)
	car, err :=  s.store.GetCarByID(ctx, id)
	if err != nil {
		return models.Car{}, err
	}
	return car, nil
}

func (s *CarService) GetCarsByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error) {
	cars, err := s.store.GetCarsByBrand(ctx, brand, isEngine)
	log.Print("Brand name in car service: " ,brand)
	log.Print(cars)
	log.Print(err)
	if err != nil {
		return []models.Car{}, err
	}
	return cars, nil
}

func (s *CarService) CreateCar(ctx context.Context, car *models.CarRequest) (models.Car, error) {
	if err := models.ValidateCarRequest(*car); err != nil {
		return models.Car{}, err
	}
	createdCar, err := s.store.CreateCar(ctx, car)
	if err != nil {
		return models.Car{}, err
	}
	return createdCar, nil
}

func (s *CarService) UpdateCar(ctx context.Context, id string, car *models.CarRequest) (models.Car, error) {
	fmt.Println("id in service: ",id)
	if err := models.ValidateCarRequest(*car); err != nil {
		fmt.Println("ERROR IN VALIDATION: ",err)
		return models.Car{}, err
	}
	
	fmt.Println("id: ",id)
	updatedCar, err := s.store.UpdateCar(ctx, id, car)
	if err != nil {
		return models.Car{}, err
	}
	return updatedCar, nil
}

func (s *CarService) DeleteCar(ctx context.Context, id string) error {

	err := s.store.DeleteCar(ctx, id)
	if err != nil {
		return err
	}
	return nil
}	