package postgres

import (
	"context"
	pb "go_auth_api_gateway/genproto/auth_service"
	"go_auth_api_gateway/pkg/utils"
	"testing"
	"time"

	//"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

func createShortUrl(t *testing.T) (resp *pb.CreateShortUrlResponse) {

	repo := NewShortenerRepo(db)

	longUrl := fakeData.URL()
	hash, err := utils.GetHash([]byte(longUrl))

	assert.NoError(t, err)

	resp, err = repo.CreateShortUrl(context.Background(), &pb.CreateShortUrlRequest{
		UserId:     createUser(t).Id,
		LongUrl:    longUrl,
		ShortUrl:   hash,
		ExpireDate: time.Now().Add(time.Hour * 1).Format(time.RFC3339),
	})

	assert.NoError(t, err)
	return resp
}

func TestCreateShortUrl(t *testing.T) {

	tests := []struct {
		name    string
		give    *pb.CreateShortUrlRequest
		wantErr error
	}{
		{
			name: "SUCCESS: Create short url",
			give: &pb.CreateShortUrlRequest{
				UserId:     createUser(t).Id,
				LongUrl:    fakeData.URL(),
				ShortUrl:   fakeData.URL(),
				ExpireDate: time.Now().Add(time.Hour * 1).Format(time.RFC3339),
			},
			wantErr: nil,
		},
		{
			name: "ERROR: Create short url, User not found",
			give: &pb.CreateShortUrlRequest{
				UserId:     createUser(t).Id,
				LongUrl:    fakeData.URL(),
				ShortUrl:   fakeData.URL(),
				ExpireDate: time.Now().Add(time.Hour * 1).Format(time.RFC3339),
			},
			wantErr: pgx.ErrNoRows,
		},
	}

	Repo := NewShortenerRepo(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Repo.CreateShortUrl(context.Background(), tt.give)
			assert.NoError(t, err)
		})
	}
}

func TestIncClickCount(t *testing.T) {

	tests := []struct {
		name    string
		give    *pb.IncClickCountRequest
		wantErr error
	}{
		{
			name: "SUCCESS: ShortUrl click count",
			give: &pb.IncClickCountRequest{
				ShortUrl: createShortUrl(t).ShortUrl,
			},
			wantErr: nil,
		},
		{
			name: "ERROR: Count Click, User not found",
			give: &pb.IncClickCountRequest{
				ShortUrl: createShortUrl(t).ShortUrl,
			},
			wantErr: pgx.ErrNoRows,
		},
	}

	Repo := NewShortenerRepo(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Repo.IncClickCount(context.Background(), tt.give)
			assert.NoError(t, err)
		})
	}
}

func TestGetShortUrl(t *testing.T) {

	tests := []struct {
		name    string
		give    *pb.GetShortUrlRequest
		wantErr error
	}{
		{
			name: "SUCCESS: Get ShortUrl",
			give: &pb.GetShortUrlRequest{
				ShortUrl: createShortUrl(t).ShortUrl,
			},
			wantErr: nil,
		},
		{
			name: "ERROR: Get shorturl, Shorturl not found",
			give: &pb.GetShortUrlRequest{
				ShortUrl: createShortUrl(t).ShortUrl,
			},
			wantErr: pgx.ErrNoRows,
		},
	}

	Repo := NewShortenerRepo(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Repo.GetShortUrl(context.Background(), tt.give, false)
			assert.NoError(t, err)
		})
	}
}

func TestGetAllUserUrls(t *testing.T) {

	createShortUrl(t)

	tests := []struct {
		name    string
		give    *pb.GetAllUserUrlsRequest
		wantErr error
	}{
		{
			name: "SUCCESS: Get User all urls",
			give: &pb.GetAllUserUrlsRequest{
				UserId: createUser(t).Id,
			},
			wantErr: nil,
		},
		{
			name: "ERROR: Get User all urls, User not found",
			give: &pb.GetAllUserUrlsRequest{
				UserId: createUser(t).Id,
			},
			wantErr: pgx.ErrNoRows,
		},
		{
			name: "SUCCESS: Create all short urls",
			give: &pb.GetAllUserUrlsRequest{
				UserId: createUser(t).Id,
			},
			wantErr: pgx.ErrNoRows,
		},
	}

	Repo := NewShortenerRepo(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Repo.GetAllUserUrls(context.Background(), tt.give)
			assert.NoError(t, err)
		})
	}
}
