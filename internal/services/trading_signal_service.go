package services

import (
	"fmt"

	"github.com/omarshah0/rest-api-with-social-auth/internal/models"
	"github.com/omarshah0/rest-api-with-social-auth/internal/repositories"
)

type TradingSignalService struct {
	repo                *repositories.TradingSignalRepository
	notificationService *NotificationService
}

func NewTradingSignalService(repo *repositories.TradingSignalRepository, notificationService *NotificationService) *TradingSignalService {
	return &TradingSignalService{
		repo:                repo,
		notificationService: notificationService,
	}
}

// Create creates a new trading signal and sends notifications
func (s *TradingSignalService) Create(signal *models.TradingSignalCreate, createdBy int64) (*models.TradingSignal, error) {
	newSignal, err := s.repo.Create(signal, createdBy)
	if err != nil {
		return nil, err
	}

	// Send notifications asynchronously (don't fail if notification fails)
	go func() {
		if err := s.notificationService.SendSignalNotification(newSignal); err != nil {
			fmt.Printf("Failed to send notification for signal %d: %v\n", newSignal.ID, err)
		}
	}()

	return newSignal, nil
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

// GetSignalsForUser retrieves trading signals visible to a specific user
func (s *TradingSignalService) GetSignalsForUser(userID int64, limit, offset int) ([]models.TradingSignal, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	return s.repo.GetSignalsForUser(userID, limit, offset)
}

// CountForUser returns the total count of trading signals visible to a specific user
func (s *TradingSignalService) CountForUser(userID int64) (int64, error) {
	return s.repo.CountForUser(userID)
}
