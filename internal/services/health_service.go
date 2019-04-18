package services

import (
	"context"
	"incrementor/internal/proto"
)

// HealthCheck struct for healthz
type HealthCheck struct {
}

// NewHealthChek construct new HealthCheck
func NewHealthChek() proto.HealthServer {
	return &HealthCheck{}
}

// Check healthz
func (hc *HealthCheck) Check(ctx context.Context, r *proto.HealthCheckRequest) (*proto.HealthCheckResponse, error) {
	return &proto.HealthCheckResponse{
		Status: proto.HealthCheckResponse_SERVING,
	}, nil
}
