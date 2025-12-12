package services

import (
	"fmt"

	"github.com/omarshah0/rest-api-with-social-auth/internal/models"
	"github.com/omarshah0/rest-api-with-social-auth/internal/repositories"
)

type PackageService struct {
	repo *repositories.PackageRepository
}

func NewPackageService(repo *repositories.PackageRepository) *PackageService {
	return &PackageService{repo: repo}
}

// Create creates a new package
func (s *PackageService) Create(pkg *models.PackageCreate) (*models.Package, error) {
	return s.repo.Create(pkg)
}

// GetByID retrieves a package by ID
func (s *PackageService) GetByID(id int64) (*models.Package, error) {
	pkg, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if pkg == nil {
		return nil, fmt.Errorf("package not found")
	}
	return pkg, nil
}

// GetAll retrieves all packages
func (s *PackageService) GetAll(activeOnly bool, limit, offset int) ([]models.Package, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 100 {
		limit = 100
	}
	return s.repo.GetAll(activeOnly, limit, offset)
}

// GetByIDs retrieves multiple packages by their IDs
func (s *PackageService) GetByIDs(ids []int64) ([]models.Package, error) {
	return s.repo.GetByIDs(ids)
}

// Update updates a package
func (s *PackageService) Update(id int64, update *models.PackageUpdate) (*models.Package, error) {
	return s.repo.Update(id, update)
}

// Delete deletes a package
func (s *PackageService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// Count returns the total count of packages
func (s *PackageService) Count(activeOnly bool) (int64, error) {
	return s.repo.Count(activeOnly)
}

// CalculateTotalPrice calculates the total price for multiple packages
func (s *PackageService) CalculateTotalPrice(packageIDs []int64) (float64, error) {
	packages, err := s.repo.GetByIDs(packageIDs)
	if err != nil {
		return 0, err
	}

	if len(packages) != len(packageIDs) {
		return 0, fmt.Errorf("one or more packages not found")
	}

	var total float64
	for _, pkg := range packages {
		if !pkg.IsActive {
			return 0, fmt.Errorf("package '%s' is not active", pkg.Name)
		}
		total += pkg.Price
	}

	return total, nil
}

