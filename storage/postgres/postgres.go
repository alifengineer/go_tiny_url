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
	db *pgxpool.Pool
	// clientPlatform  storage.ClientPlatformRepoI
	// clientType      storage.ClientTypeRepoI
	// client          storage.ClientRepoI
	// role            storage.RoleRepoI
	// permission      storage.PermissionRepoI
	// scope           storage.ScopeRepoI
	// permissionScope storage.PermissionScopeRepoI
	user storage.UserRepoI
	// session         storage.SessionRepoI
	// rolePermission   storage.RolePermissionRepoI
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
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0,
	})
	redisCache := cache.New(&cache.Options{
		Redis: redisClient,
	})

	return &Store{
		db:  pool,
		rdb: redisCache,
	}, err
}

func (s *Store) CloseDB() {
	s.db.Close()
}

// func (s *Store) ClientPlatform() storage.ClientPlatformRepoI {
// 	if s.clientPlatform == nil {
// 		s.clientPlatform = NewClientPlatformRepo(s.db)
// 	}

// 	return s.clientPlatform
// }

// func (s *Store) ClientType() storage.ClientTypeRepoI {
// 	if s.clientType == nil {
// 		s.clientType = NewClientTypeRepo(s.db)
// 	}

// 	return s.clientType
// }

// func (s *Store) Client() storage.ClientRepoI {
// 	if s.client == nil {
// 		s.client = NewClientRepo(s.db)
// 	}

// 	return s.client
// }

// func (s *Store) Role() storage.RoleRepoI {
// 	if s.role == nil {
// 		s.role = NewRoleRepo(s.db)
// 	}

// 	return s.role
// }

// func (s *Store) Permission() storage.PermissionRepoI {
// 	if s.permission == nil {
// 		s.permission = NewPermissionRepo(s.db)
// 	}

// 	return s.permission
// }

// func (s *Store) Scope() storage.ScopeRepoI {
// 	if s.scope == nil {
// 		s.scope = NewScopeRepo(s.db)
// 	}

// 	return s.scope
// }

// func (s *Store) PermissionScope() storage.PermissionScopeRepoI {
// 	if s.permissionScope == nil {
// 		s.permissionScope = NewPermissionScopeRepo(s.db)
// 	}

// 	return s.permissionScope
// }

// func (s *Store) RolePermission() storage.RolePermissionRepoI {
// 	if s.rolePermission == nil {
// 		s.rolePermission = NewRolePermissionRepo(s.db)
// 	}

// 	return s.rolePermission
// }

func (s *Store) User() storage.UserRepoI {
	if s.user == nil {
		s.user = NewUserRepo(s.db)
	}

	return s.user
}

// func (s *Store) Session() storage.SessionRepoI {
// 	if s.session == nil {
// 		s.session = NewSessionRepo(s.db)
// 	}

// 	return s.session
// }

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
	