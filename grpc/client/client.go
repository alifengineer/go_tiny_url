package client

import (
	"go_auth_api_gateway/config"
	"go_auth_api_gateway/genproto/auth_service"

	"google.golang.org/grpc"
)

type ServiceManagerI interface {
	ClientService() auth_service.ClientServiceClient
	PermissionService() auth_service.PermissionServiceClient
	UserService() auth_service.UserServiceClient
	SessionService() auth_service.SessionServiceClient
}

type grpcClients struct {
	clientService     auth_service.ClientServiceClient
	permissionService auth_service.PermissionServiceClient
	userService       auth_service.UserServiceClient
	sessionService    auth_service.SessionServiceClient
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

		clientService:     auth_service.NewClientServiceClient(connAuthService),
		permissionService: auth_service.NewPermissionServiceClient(connAuthService),
		userService:       auth_service.NewUserServiceClient(connAuthService),
		sessionService:    auth_service.NewSessionServiceClient(connAuthService),
	}, nil
}

func (g *grpcClients) ClientService() auth_service.ClientServiceClient {
	return g.clientService
}

func (g *grpcClients) PermissionService() auth_service.PermissionServiceClient {
	return g.permissionService
}

func (g *grpcClients) UserService() auth_service.UserServiceClient {
	return g.userService
}

func (g *grpcClients) SessionService() auth_service.SessionServiceClient {
	return g.sessionService
}
