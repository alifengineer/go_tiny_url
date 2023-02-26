package service

import (
	"context"
	"go_auth_api_gateway/config"
	pb "go_auth_api_gateway/genproto/auth_service"
	"go_auth_api_gateway/storage"

	"github.com/saidamir98/udevs_pkg/logger"
)

type GetUrlService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	pb.GetShortUrlResponse
}

func NewGEtService(cfg config.Config, log logger.LoggerI, strg storage.StorageI) *GetUrlService{
	return &GetUrlService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
	}
}

func (s *GetUrlService) GetUrlByID(ctx context.Context, req *pb.GetShortUrlResponse) (interface{}, error) {
	s.log.Info("---GetUrlByID--->", logger.Any("req", req))

	res, err := s.strg.RedisRepo().Get(req.ShortUrl)

	if err != nil {
		s.log.Error("!!!GetUrlByID--->", logger.Error(err))
	}

	return res, nil
}


func (s *GetUrlService) GetExists(ctx context.Context, req *pb.GetShortUrlResponse) (interface{}, error) {
	s.log.Info("---GetExists--->", logger.Any("req", req))

	res, err := s.strg.RedisRepo().Exists(req.ShortUrl)

	if err != nil {
		s.log.Error("!!!GetUrlByID--->", logger.Error(err))
	}

	return res, nil
}

func (s *GetUrlService) GetSet(ctx context.Context, req *pb.GetShortUrlResponse) (error) {
	s.log.Info("---GetExists--->", logger.Any("req", req))
	
	 err := s.strg.RedisRepo().Set(req.ShortUrl, req.LongUrl)

	if err != nil {
		s.log.Error("!!!GetUrlByID--->", logger.Error(err))
	}

	return nil
}
