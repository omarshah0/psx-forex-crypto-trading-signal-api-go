package services

import (
	"fmt"

	"github.com/omarshah0/rest-api-with-social-auth/internal/models"
	"github.com/omarshah0/rest-api-with-social-auth/internal/repositories"
)

type TradingSignalService struct {
	repo *repositories.TradingSignalRepository
}

func NewTradingSignalService(repo *repositories.TradingSignalRepository) *TradingSignalService {
	return &TradingSignalService{repo: repo}
}

// Create creates a new trading signal
func (s *TradingSignalService) Create(signal *models.TradingSignalCreate, createdBy int64) (*models.TradingSignal, error) {
	return s.repo.Create(signal, createdBy)
}

// GetByID retrieves a trading signal by ID
func (s *TradingSignalService) GetByID(id int64) (*models.TradingSignal, error) {
	signal, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if signal == nil {
		return nil, fmt.Errorf("trading signal not found")
	}
	return signal, nil
}

// GetAll retrieves all trading signals
func (s *TradingSignalService) GetAll(limit, offset int) ([]models.TradingSignal, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	return s.repo.GetAll(limit, offset)
}

// Update updates a trading signal
func (s *TradingSignalService) Update(id int64, update *models.TradingSignalUpdate) (*models.TradingSignal, error) {
	return s.repo.Update(id, update)
}

// Delete deletes a trading signal
func (s *TradingSignalService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// Count returns the total count of trading signals
func (s *TradingSignalService) Count() (int64, error) {
	return s.repo.Count()
}
