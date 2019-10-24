package services

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/twinj/uuid"
	"incrementor/internal/auth"
	"incrementor/internal/interfaces"
	"incrementor/internal/models"
	"incrementor/internal/proto"
)

// IncrementorService struct hold info about Incrementor service
type IncrementorService struct {
	IncRepo interfaces.IncrementRepository

	ClientSrv *ClientService
	JWT       *auth.Jwt
}

// NewIncrementorService consturct NewIncrementorService
func NewIncrementorService(incRepo interfaces.IncrementRepository, jwt *auth.Jwt, clientSrvc *ClientService) proto.IncrementorServiceServer {
	return &IncrementorService{
		IncRepo:   incRepo,
		JWT:       jwt,
		ClientSrv: clientSrvc,
	}
}

func (s *IncrementorService) clientIncremen(ctx context.Context) (*models.Increment, error) {
	username := ctx.Value(auth.ContextKey).(string)
	inc, err := s.IncRepo.FindByUsername(username)
	if err != nil {
		return nil, errors.Wrap(err, "error find client incrementor")
	}
	return inc, nil
}

// GetNumber return current number
func (s *IncrementorService) GetNumber(ctx context.Context, e *empty.Empty) (*proto.GetNumberResponse, error) {

	increment, err := s.clientIncremen(ctx)
	if err != nil {
		return nil, err
	}

	return &proto.GetNumberResponse{
		Number: increment.Number,
	}, nil
}

// IncrementNumber increment current number
func (s *IncrementorService) IncrementNumber(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {

	increment, err := s.clientIncremen(ctx)
	if err != nil {
		return nil, err
	}

	newNumber := increment.Number + increment.Step
	if newNumber > increment.MaxValue {
		increment.Number = 0
	} else {
		increment.Number = newNumber
	}

	increment, err = s.IncRepo.Update(increment)
	if err != nil {
		return nil, errors.Wrap(err, "Error on store increment")
	}

	return &empty.Empty{}, nil
}

// SetSettings setup settings for incrementor
func (s *IncrementorService) SetSettings(ctx context.Context, r *proto.SetSettingsRequest) (*empty.Empty, error) {

	increment, err := s.clientIncremen(ctx)
	if err != nil {
		return nil, err
	}

	if r.MaxValue > 0 {
		increment.MaxValue = r.MaxValue
	} else {
		return nil, errors.New("MaxValue must be greater than zero")
	}

	if r.IncrementSize > 0 {
		increment.Step = r.IncrementSize
	} else {
		return nil, errors.New("IncerementSize must be greater than zero")
	}

	increment, err = s.IncRepo.Update(increment)
	if err != nil {
		return nil, errors.Wrap(err, "Error on store increment")
	}

	return &empty.Empty{}, nil
}

// Auth authenticate client on incrementor service
func (s *IncrementorService) Auth(ctx context.Context, req *proto.AuthRequest) (*proto.AuthResponse, error) {

	client, err := s.ClientSrv.Auth(req.Username, req.Password)
	if err != nil {
		return nil, errors.Wrap(err, "Auth required")
	}

	tokenString, err := s.JWT.Sign(client)
	if err != nil {
		return nil, errors.Wrap(err, "Error signing jwt token")
	}

	return &proto.AuthResponse{
		Token: tokenString,
	}, nil
}

// Register new client for incrementor
func (s *IncrementorService) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	client, err := s.ClientSrv.Register(req.Username, req.Password)
	if err != nil {
		return nil, errors.Wrap(err, "Error serving RegisterRequest")
	}

	inc := &models.Increment{
		ClientGUID: client.GUID,
		Number:     models.MinValue,
		MaxValue:   models.MaxValue,
		Step:       models.IncremetStep,
	}
	inc.GUID = uuid.NewV4().String()

	inc, err = s.IncRepo.Create(inc)
	if err != nil {
		return nil, errors.Wrap(err, "Error serving RegisterRequest")
	}

	return &proto.RegisterResponse{
		GUID: client.GUID,
	}, nil

}
