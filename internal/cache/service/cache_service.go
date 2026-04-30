package cache_service

import (
	"BAZ/internal/cache"
	"context"
	"time"
)

type CacheService[T any] struct {
	cache cache.Cache
}

func NewCacheService[T any](cache cache.Cache) *CacheService[T] {
	return &CacheService[T]{cache: cache}
}

func (s *CacheService[T]) GetOrSet(
	ctx context.Context,
	key string,
	ttl time.Duration,
	fetcher func() (T, error),
) (T, error) {
	var result T

	err := s.cache.GetModel(ctx, key, &result)
	if err == nil {
		return result, nil
	}

	result, err = fetcher()
	if err != nil {
		return result, err
	}

	_ = s.cache.SetModel(ctx, key, result, ttl)

	return result, nil
}
