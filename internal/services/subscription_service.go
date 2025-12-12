package services

import (
	"fmt"
	"time"

	"github.com/omarshah0/rest-api-with-social-auth/internal/models"
	"github.com/omarshah0/rest-api-with-social-auth/internal/repositories"
)

type SubscriptionService struct {
	subscriptionRepo *repositories.SubscriptionRepository
	packageRepo      *repositories.PackageRepository
	paymentRepo      *repositories.PaymentRepository
	emailService     *EmailService
	userRepo         *repositories.UserRepository
}

func NewSubscriptionService(
	subscriptionRepo *repositories.SubscriptionRepository,
	packageRepo *repositories.PackageRepository,
	paymentRepo *repositories.PaymentRepository,
	emailService *EmailService,
	userRepo *repositories.UserRepository,
) *SubscriptionService {
	return &SubscriptionService{
		subscriptionRepo: subscriptionRepo,
		packageRepo:      packageRepo,
		paymentRepo:      paymentRepo,
		emailService:     emailService,
		userRepo:         userRepo,
	}
}

// Subscribe creates subscriptions for a user with multiple packages
func (s *SubscriptionService) Subscribe(userID int64, packageIDs []int64) (*models.SubscribeResponse, error) {
	// Fetch packages
	packages, err := s.packageRepo.GetByIDs(packageIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch packages: %w", err)
	}

	if len(packages) != len(packageIDs) {
		return nil, fmt.Errorf("one or more packages not found")
	}

	// Validate packages are active
	var totalAmount float64
	for _, pkg := range packages {
		if !pkg.IsActive {
			return nil, fmt.Errorf("package '%s' is not active", pkg.Name)
		}
		totalAmount += pkg.Price
	}

	// Create subscriptions
	var subscriptions []models.SubscriptionWithPackage
	now := time.Now()

	for i := range packages {
		pkg := packages[i] // Create a copy of the package for this iteration
		expiresAt := now.AddDate(0, 0, pkg.DurationDays)

		subscription, err := s.subscriptionRepo.Create(userID, pkg.ID, pkg.Price, expiresAt)
		if err != nil {
			return nil, fmt.Errorf("failed to create subscription: %w", err)
		}

		// Create payment record
		paymentCreate := &models.PaymentCreate{
			UserID:        userID,
			PackageID:     pkg.ID,
			Amount:        pkg.Price,
			PaymentMethod: strPtr("dummy"),
			PaymentStatus: models.PaymentStatusCompleted,
			TransactionID: strPtr(fmt.Sprintf("dummy-%d-%d-%d", userID, time.Now().UnixNano(), i)),
			Metadata:      map[string]interface{}{"note": "Dummy payment for development"},
		}

		_, err = s.paymentRepo.Create(paymentCreate)
		if err != nil {
			return nil, fmt.Errorf("failed to create payment record: %w", err)
		}

		subscriptions = append(subscriptions, models.SubscriptionWithPackage{
			Subscription: *subscription,
			Package:      &pkg,
		})
	}

	// Send confirmation email
	user, err := s.userRepo.GetByID(userID)
	if err == nil && user != nil {
		err = s.emailService.SendSubscriptionConfirmation(user.Email, user.Name, subscriptions, totalAmount)
		if err != nil {
			// Log error but don't fail the subscription
			fmt.Printf("Failed to send subscription confirmation email: %v\n", err)
		}
	}

	return &models.SubscribeResponse{
		Subscriptions: subscriptions,
		TotalAmount:   totalAmount,
		Message:       fmt.Sprintf("Successfully subscribed to %d package(s)", len(subscriptions)),
	}, nil
}

// CheckAccess checks if user has active subscription for specific asset class and duration type
func (s *SubscriptionService) CheckAccess(userID int64, assetClass models.AssetClass, durationType models.DurationType) (*models.CheckAccessResponse, error) {
	subscription, err := s.subscriptionRepo.CheckAccess(userID, assetClass, durationType)
	if err != nil {
		return nil, err
	}

	if subscription == nil {
		return &models.CheckAccessResponse{
			HasAccess: false,
		}, nil
	}

	return &models.CheckAccessResponse{
		HasAccess: true,
		ExpiresAt: &subscription.ExpiresAt,
		PackageID: &subscription.PackageID,
		PricePaid: &subscription.PricePaid,
	}, nil
}

// GetActiveSubscriptions retrieves all active subscriptions for a user
func (s *SubscriptionService) GetActiveSubscriptions(userID int64) ([]models.SubscriptionWithPackage, error) {
	subscriptions, err := s.subscriptionRepo.GetActiveByUserID(userID)
	if err != nil {
		return nil, err
	}

	var result []models.SubscriptionWithPackage
	for _, sub := range subscriptions {
		pkg, err := s.packageRepo.GetByID(sub.PackageID)
		if err != nil {
			return nil, err
		}

		result = append(result, models.SubscriptionWithPackage{
			Subscription: sub,
			Package:      pkg,
		})
	}

	return result, nil
}

// GetAllSubscriptions retrieves all subscriptions (active and expired) for a user
func (s *SubscriptionService) GetAllSubscriptions(userID int64, limit, offset int) ([]models.SubscriptionWithPackage, error) {
	subscriptions, err := s.subscriptionRepo.GetAllByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var result []models.SubscriptionWithPackage
	for _, sub := range subscriptions {
		pkg, err := s.packageRepo.GetByID(sub.PackageID)
		if err != nil {
			return nil, err
		}

		result = append(result, models.SubscriptionWithPackage{
			Subscription: sub,
			Package:      pkg,
		})
	}

	return result, nil
}

// DeactivateExpired deactivates all expired subscriptions
func (s *SubscriptionService) DeactivateExpired() (int64, error) {
	return s.subscriptionRepo.DeactivateExpired()
}

// CountByUserID returns the total count of subscriptions for a user
func (s *SubscriptionService) CountByUserID(userID int64) (int64, error) {
	return s.subscriptionRepo.CountByUserID(userID)
}

// Helper function to create string pointer
func strPtr(s string) *string {
	return &s
}
