package storage

import (
	"context"
	"errors"
	pb "go_auth_api_gateway/genproto/auth_service"
	"time"
)

var ErrorTheSameId = errors.New("cannot use the same uuid for 'id' and 'parent_id' fields")
var ErrorProjectId = errors.New("not valid 'project_id'")

type StorageI interface {
	CloseDB()
	User() UserRepoI
	Shortener() ShortenerRepoI
	RedisRepo() RedisRepoI
}

type UserRepoI interface {
	// GetListByPKs(ctx context.Context, pKeys *pb.UserPrimaryKeyList) (res *pb.GetUserListResponse, err error)
	Create(ctx context.Context, entity *pb.CreateUserRequest) (pKey *pb.UserPrimaryKey, err error)
	GetList(ctx context.Context, queryParam *pb.GetUserListRequest) (res *pb.GetUserListResponse, err error)
	GetByPK(ctx context.Context, pKey *pb.UserPrimaryKey) (res *pb.User, err error)
	Update(ctx context.Context, entity *pb.UpdateUserRequest) (rowsAffected int64, err error)
	Delete(ctx context.Context, pKey *pb.UserPrimaryKey) (rowsAffected int64, err error)
	GetByUsername(ctx context.Context, username string) (res *pb.User, err error)
	ResetPassword(ctx context.Context, user *pb.ResetPasswordRequest) (rowsAffected int64, err error)
}
type ShortenerRepoI interface {
	CreateShortUrl(ctx context.Context, req *pb.CreateShortUrlRequest) (resp *pb.CreateShortUrlResponse, err error)
	GetShortUrl(ctx context.Context, req *pb.GetShortUrlRequest) (resp *pb.GetShortUrlResponse, err error)
	IncClickCount(ctx context.Context, req *pb.IncClickCountRequest) (resp *pb.IncClickCountResponse, err error)
}

type RedisRepoI interface {
	Create(ctx context.Context, key string, obj interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string, resp interface{}) (bool, error)
	Delete(ctx context.Context, key string) error
}
