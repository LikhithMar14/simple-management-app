package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/LikhithMar14/management/models"
	"github.com/LikhithMar14/management/store"
)

type EngineService struct {
	store store.EngineStoreInterface
}

func NewEngineService(store store.EngineStoreInterface) *EngineService {
	return &EngineService{store: store}
}

func (s *EngineService) GetEngineByID(ctx context.Context, id string) (models.Engine, error) {
	engine, err := s.store.GetEngineByID(ctx, id)
	if err != nil {
		return models.Engine{}, err
	}
	return engine, nil
}

func (s *EngineService) CreateEngine(ctx context.Context, engine *models.EngineRequest) (models.Engine, error) {
	if err := models.ValidateEngineRequest(*engine); err != nil {
		return models.Engine{}, err
	}

	createdEngine, err := s.store.CreateEngine(ctx, engine)
	if err != nil {
		return models.Engine{}, err
	}
	return createdEngine, nil

}

func (s *EngineService) UpdateEngine(ctx context.Context, id string, engine *models.EngineRequest) (models.Engine, error) {
	if err := models.ValidateEngineRequest(*engine); err != nil {
		return models.Engine{}, err
	}

	updatedEngine, err := s.store.UpdateEngine(ctx, id, engine)
	if err != nil {
		return models.Engine{}, err
	}
	return updatedEngine, nil

}

func (s *EngineService) DeleteEngine(ctx context.Context, id string) error {
	engine, err := s.store.GetEngineByID(ctx, id)
	if err != nil {
		return err
	}

	if engine.EngineID == uuid.Nil {
		return errors.New("engine not found")
	}

	err = s.store.DeleteEngine(ctx, id)
	if err != nil {
		return err
	}
	return nil
}	
