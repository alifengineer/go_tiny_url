package client

import (
	"go_auth_api_gateway/config"
	"go_auth_api_gateway/genproto/auth_service"

	"google.golang.org/grpc"
)

type ServiceManagerI interface {
	UserService() auth_service.UserServiceClient
	ShortenerService() auth_service.ShortenerServiceClient
}

type grpcClients struct {
	userService      auth_service.UserServiceClient
	shortenerService auth_service.ShortenerServiceClient
}

func NewGrpcClients(cfg config.Config) (ServiceManagerI, error) {

	connAuthService, err := grpc.Dial(
		cfg.AuthServiceHost+cfg.AuthGRPCPort,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	return &grpcClients{
		userService:      auth_service.NewUserServiceClient(connAuthService),
		shortenerService: auth_service.NewShortenerServiceClient(connAuthService),
	}, nil
}

func (g *grpcClients) UserService() auth_service.UserServiceClient {
	return g.userService
}

func (g *grpcClients) ShortenerService() auth_service.ShortenerServiceClient {
	return g.shortenerService
}
