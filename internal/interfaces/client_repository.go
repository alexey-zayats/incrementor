package interfaces

import "incrementor/internal/models"

// ClientRepository describe Client manipulation interface
type ClientRepository interface {
	Create(*models.Client) (*models.Client, error)
	Update(*models.Client) (*models.Client, error)
	Delete(*models.Client) error
	FindByName(username string) (*models.Client, error)
}
