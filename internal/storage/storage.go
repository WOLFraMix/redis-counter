package storage

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Counter — то, что нужно хендлерам. Ничего больше.
type Counter interface {
	Hit(ctx context.Context, name string) (int64, error)
	Stats(ctx context.Context, name string) (int64, error)
	Top(ctx context.Context, n int64) ([]string, error)
}

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

func (s *Storage) Stats(ctx context.Context, name string) (int64, error) {
	// GET по несуществующему ключу в Redis — это nil, а не ошибка:
	// превращаем в честный 0, чтобы хендлеру не пришлось знать эту особенность
	n, err := s.rdb.Get(ctx, "hits:"+name).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return n, err
}

func (s *Storage) Hit(ctx context.Context, name string) (int64, error) {
	n, err := s.rdb.Incr(ctx, "hits:"+name).Result()
	if err != nil {
		return 0, err
	}
	// ZINCRBY: увеличивает score элемента name на 1
	if err := s.rdb.ZIncrBy(ctx, "top", 1, name).Err(); err != nil {
		return 0, err
	}
	return n, nil
}

func (s *Storage) Top(ctx context.Context, n int64) ([]string, error) {
	vals, err := s.rdb.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:   "top",
		Start: 0,
		Stop:  n - 1,
		Rev:   true, // от большего score к меньшему
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("top: %w", err)
	}
	return vals, nil
}
