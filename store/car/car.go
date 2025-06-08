package car

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

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
	log.Printf("I am in car store")

	query := `
		SELECT 
			c.id, c.name, c.year, c.brand, c.fuel_type, c.engine_id, c.price, c.created_at, c.updated_at,
			e.displacement, e.number_of_cylinders, e.car_range
		FROM cars c
		JOIN engines e ON c.engine_id = e.id
		WHERE c.id = $1
	`

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&car.ID,
		&car.Name,
		&car.Year,
		&car.Brand,
		&car.FuelType,
		&car.Engine.EngineID,
		&car.Price,
		&car.CreatedAt,
		&car.UpdatedAt,
		&car.Engine.Displacement,
		&car.Engine.NumberOfCylinders,
		&car.Engine.CarRange,
	)
	log.Println(err)

	if err != nil {
		return models.Car{}, err
	}
	return car, nil
}

func (s *CarStore) GetCarsByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error) {
	var cars []models.Car
	var query string
	log.Println(brand)
	log.Println(isEngine)

	if isEngine {
		query = `
			SELECT c.id, c.name, c.year, c.brand, c.fuel_type, c.engine_id, c.price, c.created_at, c.updated_at, 
			       e.displacement, e.number_of_cylinders, e.car_range
			FROM cars c 
			LEFT JOIN engines e ON c.engine_id = e.id
			WHERE c.brand = $1
		`
	} else {
		query = `
			SELECT c.id, c.name, c.year, c.brand, c.fuel_type, c.engine_id, c.price, c.created_at, c.updated_at
			FROM cars c
			WHERE c.brand = $1
		`
	}

	rows, err := s.db.QueryContext(ctx, query, brand)
	if err != nil {
		return []models.Car{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var car models.Car
		if isEngine {
			err := rows.Scan(
				&car.ID, &car.Name, &car.Year, &car.Brand, &car.FuelType,
				&car.Engine.EngineID, &car.Price, &car.CreatedAt, &car.UpdatedAt,
				&car.Engine.Displacement, &car.Engine.NumberOfCylinders, &car.Engine.CarRange,
			)
			if err != nil {
				return []models.Car{}, err
			}
		} else {
			err := rows.Scan(
				&car.ID, &car.Name, &car.Year, &car.Brand, &car.FuelType,
				&car.Engine.EngineID, &car.Price, &car.CreatedAt, &car.UpdatedAt,
			)
			if err != nil {
				return []models.Car{}, err
			}
		}
		cars = append(cars, car)
	}

	// Check for any error that occurred during iteration
	if err = rows.Err(); err != nil {
		return []models.Car{}, err
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

	err = tx.QueryRowContext(ctx, carQuery, car.Name, car.Year, car.Brand, car.FuelType, engineID, car.Price).Scan(
		&newCar.ID, &newCar.Name, &newCar.Year, &newCar.Brand, &newCar.FuelType, 
		&newCar.Engine.EngineID, &newCar.Price, &newCar.CreatedAt, &newCar.UpdatedAt,
	)
	if err != nil {
		return models.Car{}, err
	}


	newCar.Engine.Displacement = car.Engine.Displacement
	newCar.Engine.NumberOfCylinders = car.Engine.NumberOfCylinders
	newCar.Engine.CarRange = car.Engine.CarRange

	return newCar, nil
}

func (s *CarStore) UpdateCar(ctx context.Context, id string, car *models.CarRequest) (models.Car, error) {
	log.Println("I am in the store")
	var updatedCar models.Car

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


	if car.Engine.EngineID != uuid.Nil {
		engineUpdateQuery := `
			UPDATE engines 
			SET displacement = $1, number_of_cylinders = $2, car_range = $3
			WHERE id = $4
		`
		_, err = tx.ExecContext(ctx, engineUpdateQuery, 
			car.Engine.Displacement, car.Engine.NumberOfCylinders, car.Engine.CarRange, car.Engine.EngineID)
		if err != nil {
			return models.Car{}, err
		}
	}


	carUpdateQuery := `
		UPDATE cars
		SET name = $1, year = $2, brand = $3, fuel_type = $4, engine_id = $5, price = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
		RETURNING id, name, year, brand, fuel_type, engine_id, price, created_at, updated_at
	`

	err = tx.QueryRowContext(ctx, carUpdateQuery, 
		car.Name, car.Year, car.Brand, car.FuelType, car.Engine.EngineID, car.Price, id).Scan(
		&updatedCar.ID, &updatedCar.Name, &updatedCar.Year, &updatedCar.Brand, 
		&updatedCar.FuelType, &updatedCar.Engine.EngineID, &updatedCar.Price, 
		&updatedCar.CreatedAt, &updatedCar.UpdatedAt,
	)
	if err != nil {
		return models.Car{}, err
	}


	engineQuery := `
		SELECT displacement, number_of_cylinders, car_range
		FROM engines
		WHERE id = $1
	`
	err = tx.QueryRowContext(ctx, engineQuery, updatedCar.Engine.EngineID).Scan(
		&updatedCar.Engine.Displacement, &updatedCar.Engine.NumberOfCylinders, &updatedCar.Engine.CarRange,
	)
	if err != nil {
		return models.Car{}, err
	}
	fmt.Print("Updated Car: ",updatedCar)


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


	var engineID uuid.UUID
	getEngineQuery := `SELECT engine_id FROM cars WHERE id = $1`
	err = tx.QueryRowContext(ctx, getEngineQuery, id).Scan(&engineID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("car not found")
		}
		return err
	}


	carQuery := `DELETE FROM cars WHERE id = $1`
	result, err := tx.ExecContext(ctx, carQuery, id)
	if err != nil {
		return err
	}


	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("car not found")
	}


	engineQuery := `DELETE FROM engines WHERE id = $1`
	_, err = tx.ExecContext(ctx, engineQuery, engineID)
	if err != nil {
		return err
	}

	return nil
}