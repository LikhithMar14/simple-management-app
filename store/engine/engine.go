package engine

import (
	"context"
	"database/sql"
	"errors"
	"github.com/LikhithMar14/management/models"
)


type EngineStore struct {
	db *sql.DB
}


func NewEngineStore(db *sql.DB) *EngineStore {
	return &EngineStore{db:db}
}

func (s *EngineStore) GetEngineByID(ctx context.Context, id string) (models.Engine, error) {
	var engine models.Engine

	query := `
		SELECT id, displacement, number_of_cylinders, car_range
		FROM engines
		WHERE id = $1
	`

	err := s.db.QueryRowContext(ctx, query, id).Scan(&engine.EngineID, &engine.Displacement, &engine.NumberOfCylinders, &engine.CarRange)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Engine{}, errors.New("engine not found")
		}
		return models.Engine{}, err
	}
	return engine, nil
	}
func (s *EngineStore)CreateEngine(ctx context.Context, engine *models.EngineRequest) (models.Engine, error) {
	var newEngine models.Engine

	query := `
		INSERT INTO engines (displacement, number_of_cylinders, car_range)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	if err := models.ValidateEngineRequest(*engine); err != nil {
		return models.Engine{}, err
	}
	err := s.db.QueryRowContext(ctx, query, engine.Displacement, engine.NumberOfCylinders, engine.CarRange).Scan(&newEngine.EngineID)
	if err != nil {
		return models.Engine{}, err
	}
	return newEngine, nil
}

func (s *EngineStore)UpdateEngine(ctx context.Context, id string, engine *models.EngineRequest) (models.Engine, error) {
	if err := models.ValidateEngineRequest(*engine); err != nil {
		return models.Engine{}, err
	}
	var updatedEngine models.Engine

	query := `
		UPDATE engines
		SET displacement = $1, number_of_cylinders = $2, car_range = $3
		WHERE id = $4
	`

	err := s.db.QueryRowContext(ctx, query, engine.Displacement, engine.NumberOfCylinders, engine.CarRange, id).Scan(&updatedEngine.EngineID, &updatedEngine.Displacement, &updatedEngine.NumberOfCylinders, &updatedEngine.CarRange)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Engine{}, errors.New("engine not found")
		}
		return models.Engine{}, err
	}
	return updatedEngine, nil
}

func (s *EngineStore)DeleteEngine(ctx context.Context, id string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func () {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()	
	query := `
		DELETE FROM engines
		WHERE id = $1
	`

	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("engine not found")
		}
		return err
	}
	return nil
}