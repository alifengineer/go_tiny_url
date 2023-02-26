package postgres

import (
	"context"
	"fmt"
	"go_auth_api_gateway/config"
	"go_auth_api_gateway/storage"
	"go_auth_api_gateway/storage/redis"

	"github.com/go-redis/cache/v9"
	"github.com/jackc/pgx/v4/pgxpool"
	goRedis "github.com/redis/go-redis/v9"
)

type Store struct {
	db        *pgxpool.Pool
	user      storage.UserRepoI
	shortener storage.ShortenerRepoI
	redisRepo storage.RedisRepoI
	rdb       *cache.Cache
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

	redisClient := goRedis.NewClient(&goRedis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0,
	})
	redisCache := cache.New(&cache.Options{
		Redis: redisClient,
	})

	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		panic("Redis is not available" + err.Error())
	}

	return &Store{
		db:  pool,
		rdb: redisCache,
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

func (s *Store) RedisRepo() storage.RedisRepoI {
	if s.redisRepo == nil {
		s.redisRepo = redis.NewredisRepo(s.rdb)
	}

	return s.redisRepo
}

func (s *Store) DB() *pgxpool.Pool {
	return s.db
}
