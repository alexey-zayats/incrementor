package rpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"incrementor/internal/config"
	"incrementor/internal/proto"
	"sync"
	"time"
)

// IncrementorClient struct holds client information to coll incrementor service
type IncrementorClient struct {
	Incrementor proto.IncrementorServiceClient
	Token       string
	Username    string
	Password    string
	mu          *sync.Mutex
}

// NewIncrementorClient construct instance of IncrementorClient
func NewIncrementorClient(config *config.Config) *IncrementorClient {

	c := &IncrementorClient{
		Username: config.Client.Auth.Username,
		Password: config.Client.Auth.Password,
	}

	creds, err := credentials.NewClientTLSFromFile(config.Client.TLS.CertFile, "")
	if err != nil {
		logrus.Fatalf("Error creating client TLS: %s", err.Error())
	}

	conn, err := grpc.Dial(config.Client.Dial.Address,
		grpc.WithTransportCredentials(creds),
		grpc.WithUnaryInterceptor(c.clientInterceptor))

	if err != nil {
		logrus.Fatal("Error dial to: ", config.Client.Dial.Address)
	}

	c.Incrementor = proto.NewIncrementorServiceClient(conn)

	return c
}

func (c *IncrementorClient) register() {
	reg, err := c.Incrementor.Register(context.Background(), &proto.RegisterRequest{
		Username: c.Username,
		Password: c.Password,
	})
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	logrus.Info(reg.GUID)
}

func (c *IncrementorClient) auth() {
	auth, err := c.Incrementor.Auth(context.Background(), &proto.AuthRequest{
		Username: c.Username,
		Password: c.Password,
	})
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	c.Token = auth.Token
}

func (c *IncrementorClient) setSettings(size, max int32) {
	_, err := c.Incrementor.SetSettings(context.Background(), &proto.SetSettingsRequest{
		IncerementSize: size,
		MaxValue:       max,
	})
	if err != nil {
		logrus.Error(err.Error())
		return
	}
}

func (c *IncrementorClient) increment() {
	_, err := c.Incrementor.IncrementNumber(context.Background(), new(empty.Empty))
	if err != nil {
		logrus.Error(err.Error())
		return
	}
}

func (c *IncrementorClient) number() int32 {
	res, err := c.Incrementor.GetNumber(context.Background(), new(empty.Empty))
	if err != nil {
		logrus.Error(err.Error())
		return 0
	}
	return res.Number
}

// Run incrementorClient
func (c *IncrementorClient) Run() error {

	c.register()
	c.auth()
	c.setSettings(10, 20)
	c.increment()
	logrus.Infof("Number: %d", c.number())

	return nil
}

func (c *IncrementorClient) clientInterceptor(ctx context.Context, method string,
	req interface{}, reply interface{}, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	start := time.Now()

	if c.Token != "" {
		md := metadata.Pairs("authorization", c.Token)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	// Calls the invoker to execute RPC
	err := invoker(ctx, method, req, reply, cc, opts...)

	logrus.Infof("Invoked RPC method=%s; Duration=%s; Error=%v", method,
		time.Since(start), err)

	return err
}
