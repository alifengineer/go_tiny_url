package postgres

import (
	"context"
	"fmt"
	"go_auth_api_gateway/config"
	"go_auth_api_gateway/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	db        *pgxpool.Pool
	user      storage.UserRepoI
	shortener storage.ShortenerRepoI
}

func NewPostgres(ctx context.Context, cfg config.Config) (storage.StorageI, error) {
	config, err := pgxpool.ParseConfig(fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDatabase,
	))
	if err != nil {
		return nil, err
	}

	config.MaxConns = cfg.PostgresMaxConnections

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: pool,
	}, err
}

func (s *Store) CloseDB() {
	s.db.Close()
}

func (s *Store) User() storage.UserRepoI {
	if s.user == nil {
		s.user = NewUserRepo(s.db)
	}

	return s.user
}

func (s *Store) Shortener() storage.ShortenerRepoI {
	if s.shortener == nil {
		s.shortener = NewShortenerRepo(s.db)
	}

	return s.shortener
}
