package services

import (
	"fmt"

	"github.com/omarshah0/rest-api-with-social-auth/internal/models"
	"github.com/omarshah0/rest-api-with-social-auth/internal/repositories"
)

type PaymentService struct {
	paymentRepo *repositories.PaymentRepository
	packageRepo *repositories.PackageRepository
}

func NewPaymentService(paymentRepo *repositories.PaymentRepository, packageRepo *repositories.PackageRepository) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		packageRepo: packageRepo,
	}
}

// Create creates a new payment record
func (s *PaymentService) Create(payment *models.PaymentCreate) (*models.Payment, error) {
	// Validate package exists
	pkg, err := s.packageRepo.GetByID(payment.PackageID)
	if err != nil {
		return nil, err
	}
	if pkg == nil {
		return nil, fmt.Errorf("package not found")
	}

	return s.paymentRepo.Create(payment)
}

// GetByID retrieves a payment by ID
func (s *PaymentService) GetByID(id int64) (*models.Payment, error) {
	payment, err := s.paymentRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, fmt.Errorf("payment not found")
	}
	return payment, nil
}

// GetByUserID retrieves all payments for a user with package details
func (s *PaymentService) GetByUserID(userID int64, limit, offset int) ([]models.PaymentWithPackage, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	payments, err := s.paymentRepo.GetByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var result []models.PaymentWithPackage
	for _, payment := range payments {
		pkg, err := s.packageRepo.GetByID(payment.PackageID)
		if err != nil {
			return nil, err
		}

		result = append(result, models.PaymentWithPackage{
			Payment: payment,
			Package: pkg,
		})
	}

	return result, nil
}

// GetByTransactionID retrieves a payment by transaction ID
func (s *PaymentService) GetByTransactionID(transactionID string) (*models.Payment, error) {
	payment, err := s.paymentRepo.GetByTransactionID(transactionID)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, fmt.Errorf("payment not found")
	}
	return payment, nil
}

// UpdateStatus updates the payment status
func (s *PaymentService) UpdateStatus(id int64, status models.PaymentStatus) error {
	return s.paymentRepo.UpdateStatus(id, status)
}

// CountByUserID returns the total count of payments for a user
func (s *PaymentService) CountByUserID(userID int64) (int64, error) {
	return s.paymentRepo.CountByUserID(userID)
}

