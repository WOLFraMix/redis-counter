package storage

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Storage struct {
	rdb *redis.Client
}

func New(ctx context.Context, addr string) (*Storage, error) {
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &Storage{rdb: rdb}, nil
}
