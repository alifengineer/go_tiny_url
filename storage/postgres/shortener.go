package postgres

import (
	"context"
	pb "go_auth_api_gateway/genproto/auth_service"
	"go_auth_api_gateway/storage"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
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

	id := uuid.New().String()

	query := `
		INSERT INTO urls (
			id,
			long_url,
			short_url,
			expire_date,
			user_id,
			created_at,
			updated_at
		`

	_, err = s.db.Exec(ctx, query,
		id,
		req.GetLongUrl(),
		req.GetShortUrl(),
		req.GetExpireDate(),
		req.GetUserId(),
		time.Now().UTC().Format(time.RFC3339),
		time.Now().UTC().Format(time.RFC3339),
	)

	if err != nil {
		return nil, errors.Wrap(err, "error while inserting short url")
	}

	resp = &pb.CreateShortUrlResponse{
		UserId:     req.GetUserId(),
		LongUrl:    req.GetLongUrl(),
		ShortUrl:   req.GetShortUrl(),
		ExpireDate: req.GetExpireDate(),
	}

	return
}
func (s *shortenerRepo) GetShortUrl(ctx context.Context, req *pb.GetShortUrlRequest) (resp *pb.GetShortUrlResponse, err error) {

	query := `
		SELECT
			long_url,
			short_url,
			expire_date,
			user_id,
			created_at,
			updated_at
		FROM urls
		WHERE short_url = $1
	`

	err = s.db.QueryRow(ctx, query, req.GetShortUrl()).Scan(
		&resp.LongUrl,
		&resp.ShortUrl,
		&resp.ExpireDate,
		&resp.UserId,
		&resp.CreatedAt,
		&resp.UpdatedAt,
	)
	if err != nil {
		return nil, errors.Wrap(err, "error while getting short url")
	}

	return
}
