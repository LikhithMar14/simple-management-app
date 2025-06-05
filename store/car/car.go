package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/LikhithMar14/management/models"
	"github.com/google/uuid"
)

type CarStore struct {
	db *sql.DB
}

func NewCarStore(db *sql.DB) *CarStore {
	return &CarStore{db: db}
}

func (s *CarStore) GetCarByID(ctx context.Context, id string) (models.Car, error) {
	var car models.Car

	query := `
		SELECT c.id, c.name, c.year, c.brand, c.fuel_type, c.engine_id, c.price, c.created_at, c.updated_at
		FROM cars
		WHERE id = $1
	`

	err := s.db.QueryRowContext(ctx, query, id).Scan(&car.ID, &car.Name, &car.Year, &car.Brand, &car.FuelType, &car.Engine.EngineID, &car.Price, &car.CreatedAt, &car.UpdatedAt, &car.Engine.Displacement, &car.Engine.NumberOfCylinders, &car.Engine.CarRange)
	if err != nil {
		return models.Car{}, err
	}
	return car, nil
}

func (s *CarStore) GetCarsByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error) {
	var cars []models.Car
	var query string

	if isEngine {
		query = `
			SELECT c.id, c.name, c.year, c.brand, c.fuel_type, c.engine_id, c.price, c.created_at, c.updated_at, e.displacement, e.number_of_cylinders, e.car_range
			FROM cars LEFT JOIN engines e ON c.engine_id = e.id
			WHERE brand = $1
		`
	} else {
		query = `
			SELECT c.id, c.name, c.year, c.brand, c.fuel_type, c.engine_id, c.price, c.created_at, c.updated_at
			FROM cars
			WHERE brand = $1
		`
	}

	rows, err := s.db.QueryContext(ctx, query, brand)

	if err != nil {
		return []models.Car{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var car models.Car
		err := rows.Scan(&car.ID, &car.Name, &car.Year, &car.Brand, &car.FuelType, &car.Engine.EngineID, &car.Price, &car.CreatedAt, &car.UpdatedAt, &car.Engine.Displacement, &car.Engine.NumberOfCylinders, &car.Engine.CarRange)
		if err != nil {
			return []models.Car{}, err
		}
		cars = append(cars, car)
	}

	return cars, nil
}

func (s *CarStore) CreateCar(ctx context.Context, car *models.CarRequest) (models.Car, error) {
	var newCar models.Car

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Car{}, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	engineQuery := `
		INSERT INTO engines (displacement, number_of_cylinders, car_range)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var engineID uuid.UUID
	err = tx.QueryRowContext(ctx, engineQuery, car.Engine.Displacement, car.Engine.NumberOfCylinders, car.Engine.CarRange).Scan(&engineID)
	if err != nil {
		return models.Car{}, err
	}

	carQuery := `
		INSERT INTO cars (name, year, brand, fuel_type, engine_id, price)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, year, brand, fuel_type, engine_id, price, created_at, updated_at
	`

	err = tx.QueryRowContext(ctx, carQuery, car.Name, car.Year, car.Brand, car.FuelType, engineID, car.Price).Scan(&newCar.ID, &newCar.Name, &newCar.Year, &newCar.Brand, &newCar.FuelType, &newCar.Engine.EngineID, &newCar.Price, &newCar.CreatedAt, &newCar.UpdatedAt)
	if err != nil {
		return models.Car{}, err
	}

	return newCar, nil
}

func (s *CarStore) UpdateCar(ctx context.Context, id string, car *models.CarRequest) (models.Car, error) {
	var updatedCar models.Car

	query := `
		UPDATE cars
		SET name = $1, year = $2, brand = $3, fuel_type = $4, engine_id = $5, price = $6
		WHERE id = $7
	`

	err := s.db.QueryRowContext(ctx, query, car.Name, car.Year, car.Brand, car.FuelType, car.Engine.EngineID, car.Price, id).Scan(&updatedCar.ID, &updatedCar.Name, &updatedCar.Year, &updatedCar.Brand, &updatedCar.FuelType, &updatedCar.Engine.EngineID, &updatedCar.Price, &updatedCar.CreatedAt, &updatedCar.UpdatedAt)
	if err != nil {
		return models.Car{}, err
	}
	return updatedCar, nil
}

func (s *CarStore) DeleteCar(ctx context.Context, id string) error {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	carQuery := `
		DELETE FROM cars
		WHERE id = $1
`
	_, err = tx.ExecContext(ctx, carQuery, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("car not found")
		}
		return err	
	}
	return nil	
}
