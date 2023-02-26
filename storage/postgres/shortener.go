package postgres

import (
	"context"
	"database/sql"
	pb "go_auth_api_gateway/genproto/auth_service"
	"go_auth_api_gateway/pkg/helper"
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

	if req.GetExpireDate() == "" {
		req.ExpireDate = time.Now().Add(time.Hour * 1000).Format(time.RFC3339)
	}

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
			updated_at,
			limit_click
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8
		)
		`

	_, err = s.db.Exec(ctx, query,
		id,
		req.GetLongUrl(),
		req.GetShortUrl(),
		req.GetExpireDate(),
		req.GetUserId(),
		time.Now().UTC().Format(time.RFC3339),
		time.Now().UTC().Format(time.RFC3339),
		req.GetLimitClick(),
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

func (s *shortenerRepo) UpdateShortUrl(ctx context.Context, req *pb.CreateShortUrlRequest) (rowsAffected int64, err error) {

	if req.GetExpireDate() == "" {
		req.ExpireDate = time.Now().Add(time.Hour * 1000).Format(time.RFC3339)
	}

	query := `
		UPDATE urls SET
			long_url = :long_url,
			short_url = :short_url,
			expire_date = :expire_date,
			updated_at = now(),
			limit_click = :limit_click
		WHERE id = :id
		`
	params := map[string]interface{}{
		"id":          req.GetId(),
		"long_url":    req.GetLongUrl(),
		"short_url":   req.GetShortUrl(),
		"expire_date": req.GetExpireDate(),
		"limit_click": req.GetLimitClick(),
	}

	q, arr := helper.ReplaceQueryParams(query, params)
	result, err := s.db.Exec(ctx, q, arr...)
	if err != nil {
		return 0, err
	}
	rowsAffected = result.RowsAffected()

	return rowsAffected, nil
}

func (s *shortenerRepo) GetShortUrl(ctx context.Context, req *pb.GetShortUrlRequest, onlyExpired bool) (resp *pb.GetShortUrlResponse, err error) {

	resp = &pb.GetShortUrlResponse{}

	var (
		expireDate sql.NullString
		createdAt  sql.NullString
		updatedAt  sql.NullString
		filter     string
		status     bool
	)

	if onlyExpired {
		filter += " AND expire_date < now() OR click_count < limit_click "
	}
	query := `
		SELECT
			long_url,
			short_url,
			expire_date,
			user_id,
			created_at,
			updated_at,
			click_count::text,
			limit_click::text,
			(CASE
				WHEN expire_date < now() OR click_count < limit_click THEN true
				ELSE false
				END) AS is_expired
		FROM urls
		WHERE short_url = $1 ` + filter + ``

	err = s.db.QueryRow(ctx, query, req.GetShortUrl()).Scan(
		&resp.LongUrl,
		&resp.ShortUrl,
		&expireDate,
		&resp.UserId,
		&createdAt,
		&updatedAt,
		&resp.ClickCount,
		&resp.LimitClick,
		&status,
	)

	if status {
		resp.UrlStatus = "Url is expired"
	} else {
		resp.UrlStatus = "Url is not expired"
	}
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
			click_count::text,
			count(1) OVER() AS total_count,
			limit_click::text,
			(CASE
				WHEN expire_date < now() OR click_count < limit_click THEN true
				ELSE false
				END) AS is_expired
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
			status     bool
		)

		url := &pb.UrlData{}

		err = rows.Scan(
			&url.Id,
			&url.LongUrl,
			&url.ShortUrl,
			&expireDate,
			&url.ClickCount,
			&resp.TotalCount,
			&url.LimitClick,
			&status,
		)
		if err != nil {
			return nil, errors.Wrap(err, "error while scanning user urls")
		}

		url.ExpireDate = expireDate.String
		if status {
			url.UrlStatus = "Url is expired"
		} else {
			url.UrlStatus = "Url is not expired"
		}

		resp.Urls = append(resp.Urls, url)
		resp.TotalCount = totalCount
	}

	return resp, nil
}
