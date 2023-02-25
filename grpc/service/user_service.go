package service

import (
	"context"
	"fmt"
	"go_auth_api_gateway/config"
	pb "go_auth_api_gateway/genproto/auth_service"
	"go_auth_api_gateway/grpc/client"
	"go_auth_api_gateway/pkg/helper"
	"go_auth_api_gateway/storage"
	"regexp"
	"time"

	"github.com/saidamir98/udevs_pkg/security"

	"github.com/saidamir98/udevs_pkg/logger"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
	pb.UnimplementedUserServiceServer
}

func NewUserService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, svcs client.ServiceManagerI) *userService {
	return &userService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: svcs,
	}
}

func (s *userService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	s.log.Info("---CreateUser--->", logger.Any("req", req))

	if len(req.Password) < 6 {
		err := fmt.Errorf("password must not be less than 6 characters")
		s.log.Error("!!!CreateUser--->", logger.Error(err))
		return nil, err
	}

	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		s.log.Error("!!!CreateUser--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	req.Password = hashedPassword

	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	email := emailRegex.MatchString(req.Email)
	if !email {
		err = fmt.Errorf("email is not valid")
		s.log.Error("!!!CreateUser--->", logger.Error(err))
		return nil, err
	}

	phoneRegex := regexp.MustCompile(`^[+]?(\d{1,2})?[\s.-]?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}$`)
	phone := phoneRegex.MatchString(req.Phone)
	if !phone {
		err = fmt.Errorf("phone number is not valid")
		s.log.Error("!!!CreateUser--->", logger.Error(err))
		return nil, err
	}

	pKey, err := s.strg.User().Create(ctx, req)

	if err != nil {
		s.log.Error("!!!CreateUser--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return s.strg.User().GetByPK(ctx, pKey)
}

func (s *userService) GetUserByID(ctx context.Context, req *pb.UserPrimaryKey) (*pb.User, error) {
	s.log.Info("---GetUserByID--->", logger.Any("req", req))

	res, err := s.strg.User().GetByPK(ctx, req)

	if err != nil {
		s.log.Error("!!!GetUserByID--->", logger.Error(err))
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return res, nil
}

func (s *userService) GetUserListByIDs(ctx context.Context, req *pb.UserPrimaryKeyList) (*pb.GetUserListResponse, error) {
	s.log.Info("---GetUserListByIDs--->", logger.Any("req", req))

	res, err := s.strg.User().GetListByPKs(ctx, req)
	if err != nil {
		s.log.Error("!!!GetUserListByIDs--->", logger.Error(err))
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return res, err
}

func (s *userService) GetUserList(ctx context.Context, req *pb.GetUserListRequest) (*pb.GetUserListResponse, error) {
	s.log.Info("---GetUserList--->", logger.Any("req", req))

	res, err := s.strg.User().GetList(ctx, req)

	if err != nil {
		s.log.Error("!!!GetUserList--->", logger.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return res, err
}

func (s *userService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	s.log.Info("---UpdateUser--->", logger.Any("req", req))

	rowsAffected, err := s.strg.User().Update(ctx, req)

	if err != nil {
		s.log.Error("!!!UpdateUser--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if rowsAffected <= 0 {
		return nil, status.Error(codes.InvalidArgument, "no rows were affected")
	}

	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	email := emailRegex.MatchString(req.Email)
	if !email {
		err = fmt.Errorf("email is not valid")
		s.log.Error("!!!UpdateUser--->", logger.Error(err))
		return nil, err
	}

	phoneRegex := regexp.MustCompile(`^[+]?(\d{1,2})?[\s.-]?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}$`)
	phone := phoneRegex.MatchString(req.Phone)
	if !phone {
		err = fmt.Errorf("phone number is not valid")
		s.log.Error("!!!UpdateUser--->", logger.Error(err))
		return nil, err
	}

	res, err := s.strg.User().GetByPK(ctx, &pb.UserPrimaryKey{Id: req.Id})
	if err != nil {
		s.log.Error("!!!UpdateUser--->", logger.Error(err))
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return res, err
}

func (s *userService) DeleteUser(ctx context.Context, req *pb.UserPrimaryKey) (*emptypb.Empty, error) {
	s.log.Info("---DeleteUser--->", logger.Any("req", req))

	res := &emptypb.Empty{}

	rowsAffected, err := s.strg.User().Delete(ctx, req)

	if err != nil {
		s.log.Error("!!!DeleteUser--->", logger.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	if rowsAffected <= 0 {
		return nil, status.Error(codes.InvalidArgument, "no rows were affected")
	}

	return res, nil
}
func (s *userService) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.User, error) {
	s.log.Info("---ResetPassword--->", logger.Any("req", req))

	if len(req.Password) < 6 {
		err := fmt.Errorf("password must not be less than 6 characters")
		s.log.Error("!!!ResetPassword--->", logger.Error(err))
		return nil, err
	}

	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		s.log.Error("!!!ResetPassword--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	tokenInfo, err := security.ParseClaims(req.Token, s.cfg.SecretKey)
	if err != nil {
		s.log.Error("!!!ResetPassword--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	req.Password = hashedPassword
	req.UserId = tokenInfo.ID

	rowsAffected, err := s.strg.User().ResetPassword(ctx, req)
	if err != nil {
		s.log.Error("!!!ResetPassword--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if rowsAffected <= 0 {
		return nil, status.Error(codes.InvalidArgument, "no rows were affected")
	}

	return s.strg.User().GetByPK(ctx, &pb.UserPrimaryKey{Id: req.UserId})
}

func (s *userService) SendMessageToEmail(ctx context.Context, req *pb.SendMessageToEmailRequest) (*emptypb.Empty, error) {
	user, err := s.strg.User().GetByUsername(context.Background(), req.GetEmail())
	if err != nil {
		s.log.Error("error while getting user by email", logger.Error(err), logger.Any("req", req))
		return nil, status.Error(codes.NotFound, err.Error())
	}

	m := map[string]interface{}{
		"id": user.Id,
	}

	token, err := security.GenerateJWT(m, time.Hour*2, s.cfg.SecretKey)
	if err != nil {
		s.log.Error("error while getting generating token", logger.Error(err), logger.Any("req", req))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = helper.SendEmail("Update Password", req.GetEmail(), req.GetBaseUrl(), token)
	if err != nil {
		s.log.Error("!!!SendUpdatePasswordUrlToEmail--->", logger.Error(err), logger.Any("req", req))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &emptypb.Empty{}, nil
}
