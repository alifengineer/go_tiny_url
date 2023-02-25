package main

import (
	"context"
	"go_auth_api_gateway/api"
	"go_auth_api_gateway/api/handlers"
	"go_auth_api_gateway/config"
	"go_auth_api_gateway/grpc"
	"go_auth_api_gateway/grpc/client"
	"go_auth_api_gateway/storage/postgres"
	"net"

	"github.com/saidamir98/udevs_pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	loggerLevel := logger.LevelDebug

	switch cfg.Environment {
	case config.DebugMode:
		loggerLevel = logger.LevelDebug
		gin.SetMode(gin.DebugMode)
	case config.TestMode:
		loggerLevel = logger.LevelDebug
		gin.SetMode(gin.TestMode)
	default:
		loggerLevel = logger.LevelInfo
		gin.SetMode(gin.ReleaseMode)
	}

	log := logger.NewLogger(cfg.ServiceName, loggerLevel)
	defer logger.Cleanup(log)

	pgStore, err := postgres.NewPostgres(context.Background(), cfg)
	if err != nil {
		log.Panic("postgres.NewPostgres", logger.Error(err))
	}
	defer pgStore.CloseDB()

	svcs, err := client.NewGrpcClients(cfg)
	if err != nil {
		log.Panic("client.NewGrpcClients", logger.Error(err))
	}

	grpcServer := grpc.SetUpServer(cfg, log, pgStore, svcs)
	go func() {
		lis, err := net.Listen("tcp", cfg.AuthGRPCPort)
		if err != nil {
			log.Panic("net.Listen", logger.Error(err))
		}

		log.Info("GRPC: Server being started...", logger.String("port", cfg.AuthGRPCPort))

		if err := grpcServer.Serve(lis); err != nil {
			log.Panic("grpcServer.Serve", logger.Error(err))
		}
	}()

	h := handlers.NewHandler(cfg, log, svcs)

	r := api.SetUpRouter(h, cfg)

	r.Run(cfg.HTTPPort)
}
