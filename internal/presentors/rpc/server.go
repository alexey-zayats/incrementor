package rpc

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"incrementor/internal/auth"
	"incrementor/internal/config"
	"incrementor/internal/proto"
	"net"
	"net/http"
	"strings"
	"time"
)

// Server - struct with info about grpc server
type Server struct {
	Config      *config.Config
	RPC         *grpc.Server
	JWT         *auth.Jwt
	SkipMethods map[string]bool
}

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

// NewServer creates new instance of rpc.Server
func NewServer(config *config.Config,
	incrementor proto.IncrementorServiceServer, j *auth.Jwt,
	healthz proto.HealthServer) (*Server, error) {

	creds, err := credentials.NewServerTLSFromFile(config.Server.TLS.CertFile, config.Server.TLS.KeyFile)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error creating ServerTLSFromFile: %s, %s", config.Server.TLS.CertFile, config.Server.TLS.KeyFile))
	}

	s := &Server{
		Config:      config,
		JWT:         j,
		SkipMethods: make(map[string]bool),
	}

	s.SkipMethods["/proto.IncrementorService/Auth"] = true
	s.SkipMethods["/proto.IncrementorService/Register"] = true

	s.RPC = grpc.NewServer(grpc.Creds(creds),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(grpc_prometheus.UnaryServerInterceptor,
			s.loggingUnary, s.monitoringUnary, s.ensureValidToken)),
	)

	proto.RegisterIncrementorServiceServer(s.RPC, incrementor)
	proto.RegisterHealthServer(s.RPC, healthz)

	grpc_prometheus.Register(s.RPC)

	http.Handle("/metrics", promhttp.Handler())

	return s, nil
}

// Run is used to start rpc server
func (s *Server) Run() error {

	lis, err := net.Listen(s.Config.Server.Listen.Network, s.Config.Server.Listen.Address)
	if err != nil {
		return errors.Wrap(err, "Error creating listener")
	}

	logrus.Infof("Serving gRPC on %s %s", s.Config.Server.Listen.Network, s.Config.Server.Listen.Address)

	if err = s.RPC.Serve(lis); err != nil {
		return errors.Wrap(err, "Error on srv.Serve(lis)")
	}

	return nil
}

// ensureValidToken ensures a valid token exists within a request's metadata
func (s *Server) ensureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	// skip Auth & Register methods
	if _, ok := s.SkipMethods[info.FullMethod]; ok {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}

	if !s.valid(md["authorization"]) {
		return nil, errInvalidToken
	}

	ctx = context.WithValue(ctx, auth.ContextKey, s.JWT.Claims.Username)

	return handler(ctx, req)
}

// valid validates the authorization.
func (s *Server) valid(authorization []string) bool {

	if len(authorization) < 1 {
		return false
	}

	token := strings.TrimPrefix(authorization[0], "Bearer ")

	result, err := s.JWT.Validate(token)
	if err != nil {
		logrus.Errorf("Jwt validation error: %s", err.Error())
	}

	return result
}

func (s *Server) loggingUnary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	//logrus.Infof("loggingUnary: %#v", info)
	return handler(ctx, req)
}

func (s *Server) monitoringUnary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	//logrus.Infof("monitoringUnary: %#v", info)

	start := time.Now()

	iface, err := handler(ctx, req)

	logrus.Infof("Invoked RPC method=%s; Duration=%s; Error=%v", info.FullMethod, time.Since(start), err)

	return iface, err
}
