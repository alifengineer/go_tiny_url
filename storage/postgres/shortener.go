package postgres

import (
	"context"
	"database/sql"
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

	var (
		expireDate sql.NullString
	)

	id := uuid.New().String()

	resp = &pb.CreateShortUrlResponse{}

	query := `
		INSERT INTO urls (
			id,
			long_url,
			short_url,
			expire_date,
			user_id,
			created_at,
			updated_at
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7
		)
		`

	_, err = s.db.Exec(ctx, query,
		id,
		req.GetLongUrl(),
		req.GetShortUrl(),
		time.Now().Add(time.Hour*1).Format(time.RFC3339),
		req.GetUserId(),
		time.Now().UTC().Format(time.RFC3339),
		time.Now().UTC().Format(time.RFC3339),
	)

	if err != nil {
		return nil, errors.Wrap(err, "error while inserting short url")
	}

	err = s.db.QueryRow(
		ctx,
		`SELECT 
			id,
			long_url,
			short_url,
			expire_date		 
		FROM urls WHERE id = $1`,
		id,
	).Scan(
		&resp.Id,
		&resp.LongUrl,
		&resp.ShortUrl,
		&expireDate,
	)
	if err != nil {
		return nil, errors.Wrap(err, "error while getting short url")
	}
	resp.ExpireDate = expireDate.String

	return resp, nil
}
func (s *shortenerRepo) GetShortUrl(ctx context.Context, req *pb.GetShortUrlRequest) (resp *pb.GetShortUrlResponse, err error) {

	resp = &pb.GetShortUrlResponse{}

	var (
		expireDate sql.NullString
		createdAt  sql.NullString
		updatedAt  sql.NullString
	)

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
		&expireDate,
		&resp.UserId,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, errors.Wrap(err, "error while getting short url")
	}
	resp.CreatedAt = createdAt.String
	resp.UpdatedAt = updatedAt.String
	resp.ExpireDate = expireDate.String

	return
}

func (s *shortenerRepo) IncClickCount(ctx context.Context, req *pb.IncClickCountRequest) (resp *pb.IncClickCountResponse, err error) {

	resp = &pb.IncClickCountResponse{}

	query := `
		UPDATE urls
		SET click_count = click_count + 1
		WHERE short_url = $1
	`

	_, err = s.db.Exec(ctx, query, req.GetShortUrl())
	if err != nil {
		return nil, errors.Wrap(err, "error while incrementing click count")
	}

	err = s.db.QueryRow(ctx, `SELECT click_count FROM urls WHERE short_url = $1`, req.GetShortUrl()).Scan(&resp.ClickCount)
	if err != nil {
		return nil, errors.Wrap(err, "error while getting click count")
	}

	return resp, nil
}

func (s *shortenerRepo) GetAllUserUrls(ctx context.Context, req *pb.GetAllUserUrlsRequest) (resp *pb.GetAllUserUrlsResponse, err error) {

	resp = &pb.GetAllUserUrlsResponse{}
	var (
		totalCount int64
	)

	query := `
		SELECT
			id,
			long_url,
			short_url,
			expire_date,
			click_count,
			count(1) OVER() AS total_count
		FROM urls
		WHERE user_id = $1
	`

	rows, err := s.db.Query(ctx, query, req.GetUserId())
	if err != nil {
		return nil, errors.Wrap(err, "error while getting all user urls")
	}

	defer rows.Close()

	for rows.Next() {
		var (
			expireDate sql.NullString
		)

		url := &pb.UrlData{}

		err = rows.Scan(
			&url.Id,
			&url.LongUrl,
			&url.ShortUrl,
			&expireDate,
			&url.ClickCount,
		)
		if err != nil {
			return nil, errors.Wrap(err, "error while scanning user urls")
		}

		url.ExpireDate = expireDate.String

		resp.Urls = append(resp.Urls, url)
		resp.TotalCount = totalCount
	}

	return resp, nil
}
