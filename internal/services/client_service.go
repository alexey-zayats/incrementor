package services

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/twinj/uuid"

	"incrementor/internal/interfaces"
	"incrementor/internal/models"
)

// ClientService struct for client service
type ClientService struct {
	repo interfaces.ClientRepository
}

// NewClientService construct ClientSrv
func NewClientService(repo interfaces.ClientRepository) *ClientService {
	return &ClientService{
		repo: repo,
	}
}

// Auth authenticate client request
func (s *ClientService) Auth(user, pass string) (*models.Client, error) {
	client, err := s.repo.FindByName(user)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Client not found")
	}

	if client.Password != pass {
		return nil, status.Error(codes.Unauthenticated, "Password mismatch")
	}

	return client, nil
}

// Register client with given username and password. Returns error if clinet already exists
func (s *ClientService) Register(username, password string) (*models.Client, error) {

	client, _ := s.repo.FindByName(username)
	if client != nil {
		return nil, errors.Errorf("client already exists with username %s", username)
	}

	client = &models.Client{
		Username: username,
		Password: password,
	}
	client.GUID = uuid.NewV4().String()

	client, err := s.repo.Create(client)
	if err != nil {
		return nil, errors.Wrap(err, "error on registering client")
	}

	return client, nil
}
