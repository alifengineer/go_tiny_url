package service

import (
	"context"
	"go_auth_api_gateway/config"
	pb "go_auth_api_gateway/genproto/auth_service"
	"go_auth_api_gateway/grpc/client"
	"go_auth_api_gateway/storage"

	"github.com/saidamir98/udevs_pkg/logger"
)

type shortenerService struct {
	cfg     config.Config
	log     logger.LoggerI
	strg    storage.StorageI
	service client.ServiceManagerI
	pb.UnimplementedShortenerServiceServer
}

func NewShortenerService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, svcs client.ServiceManagerI) *shortenerService {
	return &shortenerService{
		cfg:     cfg,
		log:     log,
		strg:    strg,
		service: svcs,
	}
}

func (self *shortenerService) CreateShortUrl(ctx context.Context, req *pb.CreateShortUrlRequest) (resp *pb.CreateShortUrlResponse, err error) {

	

	return
}
func (self *shortenerService) GetShortUrl(ctx context.Context, req *pb.GetShortUrlRequest) (resp *pb.GetShortUrlResponse, err error) {

	return
}
