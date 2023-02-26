package storage

import (
	"context"
	"errors"
	pb "go_auth_api_gateway/genproto/auth_service"
)

var ErrorTheSameId = errors.New("cannot use the same uuid for 'id' and 'parent_id' fields")
var ErrorProjectId = errors.New("not valid 'project_id'")

type StorageI interface {
	CloseDB()
	User() UserRepoI
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
