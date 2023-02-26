package postgres

import (
	"context"
	pb "go_auth_api_gateway/genproto/auth_service"
	"go_auth_api_gateway/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

type shortenerRepo struct {
	db *pgxpool.Pool
}

func NewShortenerRepo(db *pgxpool.Pool) storage.ShortenerRepoI {
	return &shortenerRepo{
		db: db,
	}
}

func (s *shortenerRepo) CreateShortUrl(ctx context.Context, req *pb.CreateShortUrlRequest) (resp *pb.CreateShortUrlResponse, err error) {

	query := `
		INSERT INTO "shortener" (
			url,

	return
}
func (s *shortenerRepo) GetShortUrl(ctx context.Context, req *pb.GetShortUrlRequest) (resp *pb.GetShortUrlResponse, err error) {

	return
}
