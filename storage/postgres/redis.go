package postgres

import (
	repo "go_auth_api_gateway/storage"

	"github.com/gomodule/redigo/redis"
)

type redisRepo struct {
	Rds *redis.Pool
}

func NewRedisRepo(rds *redis.Pool) repo.RedisRepoI {
	return &redisRepo{
		Rds: rds,
	}
}

func (r *redisRepo) Exists(key string) (interface{}, error) {
	conn := r.Rds.Get()
	defer conn.Close()
	return conn.Do("EXISTS", key) 
}

func (r *redisRepo) Set(key string, value string) error {
	conn := r.Rds.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)

	return err
}

func (r *redisRepo) Get(key string) (interface{}, error) {
	conn := r.Rds.Get()
	defer conn.Close()

	return conn.Do("GET", key)
}
