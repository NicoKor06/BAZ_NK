package cache

import (
	"context"
	"errors"
	"time"
)

type MockCache struct {
	Data map[string]string
}

func NewMockCache() *MockCache {
	return &MockCache{
		Data: make(map[string]string),
	}
}

func (m *MockCache) Get(ctx context.Context, key string) (string, error) {
	val, exists := m.Data[key]
	if !exists {
		return "", nil
	}
	return val, nil
}

func (m *MockCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("mock cache currently only supports string values")
	}

	m.Data[key] = str
	return nil
}

func (m *MockCache) Delete(ctx context.Context, key string) error {
	delete(m.Data, key)
	return nil
}

func (m *MockCache) Exists(ctx context.Context, key string) (bool, error) {
	_, exists := m.Data[key]
	return exists, nil
}

func (m *MockCache) Increment(ctx context.Context, key string) (int64, error) {
	return 0, nil
}

func (m *MockCache) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	return 0, nil
}

func (m *MockCache) FlushAll(ctx context.Context) error {
	m.Data = make(map[string]string)
	return nil
}
