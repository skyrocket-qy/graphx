package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/skyrocketOoO/zanazibar-dag/domain"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) (*RedisRepository, error) {
	return &RedisRepository{
		client: client,
	}, nil
}

func (r *RedisRepository) Ping(c context.Context) error {
	return r.client.Ping(c).Err()
}

func (r *RedisRepository) Get(c context.Context, edge domain.Edge,
	queryMode bool) ([]domain.Edge, error) {
	return nil, domain.ErrNotImplemented{}
}

func (r *RedisRepository) Create(c context.Context, edge domain.Edge) error {
	return domain.ErrNotImplemented{}
}

func (r *RedisRepository) Delete(c context.Context, edge domain.Edge,
	queryMode bool) error {
	return domain.ErrNotImplemented{}
}

func (r *RedisRepository) ClearAll(c context.Context) error {
	return domain.ErrNotImplemented{}
}
