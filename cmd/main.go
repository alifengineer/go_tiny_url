package main

import (
	"context"
	"database/sql"
	"fmt"
	"go_auth_api_gateway/api"
	"go_auth_api_gateway/api/handlers"
	"go_auth_api_gateway/config"
	"go_auth_api_gateway/grpc"
	"go_auth_api_gateway/grpc/client"
	"go_auth_api_gateway/storage/postgres"
	"net"

	"github.com/saidamir98/udevs_pkg/logger"

	"github.com/gin-gonic/gin"

	migrate "github.com/golang-migrate/migrate/v4"
	pm "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	db, err := sql.Open("postgres", fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDatabase,
	))
	if err != nil {
		log.Fatal(err.Error())
	}

	driver, err := pm.WithInstance(db, &pm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://app/migrations/postgres",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err.Error())
	}

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

	h := handlers.NewHandler(cfg, log, svcs, pgStore)

	r := api.SetUpRouter(h, cfg)

	r.Run(cfg.HTTPPort)
}
