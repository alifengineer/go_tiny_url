package redis

import (
	"context"
	"go_auth_api_gateway/storage"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/pkg/errors"
)

type redisRepo struct {
	cache *cache.Cache
}

func NewredisRepo(cache *cache.Cache) storage.RedisRepoI {
	return &redisRepo{cache: cache}
}

func (u redisRepo) Create(ctx context.Context, key string, obj interface{}, ttl time.Duration) error {
	err := u.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: obj,
		TTL:   ttl,
	})
	if err != nil {
		//println("redis.Create.Error:", err.Error(), "\nkey:", id)
		return errors.Wrap(err, "error while creating cache in redis")
	}
	//fmt.Println("create in redis", id)
	return nil
}

func (u redisRepo) Get(ctx context.Context, key string, response interface{}) (bool, error) {
	// var response interface{}

	err := u.cache.Get(ctx, key, response)
	if err != nil {
		//println("redis.Get.Error:", err.Error(), "\nkey:", id)
		return false, errors.Wrap(err, "error while getting cache in redis")
	}
	//fmt.Println("get from redis", id)
	return true, nil
}

func (u redisRepo) Delete(ctx context.Context, id string) error {
	err := u.cache.Delete(ctx, id)
	if err != nil {
		return errors.Wrap(err, "error while deleting cache in redis")
	}
	//fmt.Println("delete from redis", id)
	return nil
}
