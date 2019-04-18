package interfaces

import "incrementor/internal/models"

// IncrementRepository describe Increment manipulation interface
type IncrementRepository interface {
	Create(*models.Increment) (*models.Increment, error)
	Update(*models.Increment) (*models.Increment, error)
	Delete(*models.Increment) error
	FindByUsername(string) (*models.Increment, error)
}
