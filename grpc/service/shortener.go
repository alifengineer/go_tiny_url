package service

import (
	"context"
	"fmt"
	"go_auth_api_gateway/config"
	pb "go_auth_api_gateway/genproto/auth_service"
	"go_auth_api_gateway/grpc/client"
	"go_auth_api_gateway/pkg/utils"
	"go_auth_api_gateway/storage"

	"github.com/saidamir98/udevs_pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *shortenerService) CreateShortUrl(ctx context.Context, req *pb.CreateShortUrlRequest) (resp *pb.CreateShortUrlResponse, err error) {

	s.log.Info("---CreateShortUrl--->", logger.Any("req", req))

	if !utils.IsLongCorrect(string(req.GetLongUrl())) {
		err = fmt.Errorf(fmt.Sprintf(utils.InvalidURLError, req.GetLongUrl()))

		s.log.Error("!!!CreateShortUrl--->", logger.Error(err))

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	hash, err := utils.GetHash([]byte(req.GetLongUrl()))
	if err != nil {
		s.log.Error("!!!CreateShortUrl--->", logger.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	req.ShortUrl = hash

	resp, err = s.strg.Shortener().CreateShortUrl(ctx, req)
	if err != nil {
		s.log.Error("!!!CreateShortUrl--->", logger.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = s.strg.RedisRepo().Create(ctx, hash, req.GetLongUrl(), config.RedisCacheTTL)
	if err != nil {
		s.log.Error("!!!CreateShortUrl--->", logger.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return
}

func (s *shortenerService) GetShortUrl(ctx context.Context, req *pb.GetShortUrlRequest) (resp *pb.GetShortUrlResponse, err error) {

	s.log.Info("---GetShortUrl--->", logger.Any("req", req))

	ok, err := s.strg.RedisRepo().Get(ctx, req.GetShortUrl(), resp)
	if err != nil {
		s.log.Error("!!!CreateShortUrl--->", logger.Error(err))
	}

	if !ok {
		resp, err = s.strg.Shortener().GetShortUrl(ctx, req)
		if err != nil {
			s.log.Error("!!!GetShortUrl--->", logger.Error(err))
			return nil, status.Error(codes.Internal, err.Error())
		}

		err = s.strg.RedisRepo().Create(ctx, req.GetShortUrl(), resp.GetLongUrl(), config.RedisCacheTTL)
	}

	return resp, err
}

func (s *shortenerService) IncClickCount(ctx context.Context, req *pb.IncClickCountRequest) (resp *pb.IncClickCountResponse, err error) {

	s.log.Info("---IncClickCount--->", logger.Any("req", req))

	resp, err = s.strg.Shortener().IncClickCount(ctx, req)
	if err != nil {
		s.log.Error("!!!IncClickCount--->", logger.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return
}

func (s *shortenerService) HandleLongUrl(ctx context.Context, req *pb.HandleLongUrlRequest) (resp *pb.HandleLongUrlResponse, err error) {

	s.log.Info("---HandlerLongUrl--->", logger.Any("req", req))

	var (
		longUrl string
	)

	ok, err := s.strg.RedisRepo().Get(ctx, req.GetShortUrl(), longUrl)
	if err != nil {
		s.log.Error("!!!CreateShortUrl--->", logger.Error(err))
	}

	if !ok {
		respShortUrl, err := s.strg.Shortener().GetShortUrl(ctx, &pb.GetShortUrlRequest{ShortUrl: req.GetShortUrl()})
		if err != nil {
			s.log.Error("!!!HandlerLongUrl--->", logger.Error(err))
			return nil, status.Error(codes.Internal, err.Error())
		}
		longUrl = respShortUrl.GetLongUrl()

		err = s.strg.RedisRepo().Create(ctx, req.GetShortUrl(), longUrl, config.RedisCacheTTL)
		if err != nil {
			s.log.Error("!!!HandlerLongUrl--->", logger.Error(err))
		}
	}

	resp.LongUrl = longUrl

	return
}
