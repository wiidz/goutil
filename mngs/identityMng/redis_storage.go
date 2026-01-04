package identityMng

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"
	"sync"
	"time"

	"github.com/click33/sa-token-go/core/adapter"
	"github.com/click33/sa-token-go/core/security"
	"github.com/go-redis/redis/v9"
)

// RedisStorage implements Sa-Token storage on Redis.
type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
	debug  bool
}

func init() {
	var once sync.Once
	once.Do(func() {
		// Register common types that Sa-Token stores.
		gob.Register(&security.RefreshTokenInfo{})
		gob.Register(map[string]any{})
		gob.Register([]string{})
		gob.Register(int64(0))
		gob.Register("")
	})
}

// NewRedisStorage creates a Storage backed by go-redis.
func NewRedisStorage(client *redis.Client, debug bool) adapter.Storage {
	return &RedisStorage{client: client, ctx: context.Background(), debug: debug}
}

func (s *RedisStorage) Set(key string, value any, expiration time.Duration) error {
	s.dbg("Set key=%s value=%v expiration=%v", key, value, expiration)
	if s.client == nil {
		return redis.ErrClosed
	}
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		return err
	}
	return s.client.Set(s.ctx, key, buf.Bytes(), expiration).Err()
}

func (s *RedisStorage) Get(key string) (any, error) {
	s.dbg("Get key=%s", key)
	if s.client == nil {
		return nil, redis.ErrClosed
	}
	b, err := s.client.Get(s.ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var out any
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *RedisStorage) Delete(keys ...string) error {
	s.dbg("Delete keys=%v", keys)
	if s.client == nil {
		return redis.ErrClosed
	}
	return s.client.Del(s.ctx, keys...).Err()
}

func (s *RedisStorage) Exists(key string) bool {
	s.dbg("Exists key=%s", key)
	if s.client == nil {
		return false
	}
	n, err := s.client.Exists(s.ctx, key).Result()
	return err == nil && n > 0
}

func (s *RedisStorage) Keys(pattern string) ([]string, error) {
	if s.client == nil {
		return nil, redis.ErrClosed
	}
	return s.client.Keys(s.ctx, pattern).Result()
}

func (s *RedisStorage) Expire(key string, expiration time.Duration) error {
	if s.client == nil {
		return redis.ErrClosed
	}
	return s.client.Expire(s.ctx, key, expiration).Err()
}

func (s *RedisStorage) TTL(key string) (time.Duration, error) {
	if s.client == nil {
		return 0, redis.ErrClosed
	}
	return s.client.TTL(s.ctx, key).Result()
}

func (s *RedisStorage) Clear() error {
	if s.client == nil {
		return redis.ErrClosed
	}
	return s.client.FlushDB(s.ctx).Err()
}

func (s *RedisStorage) Ping() error {
	if s.client == nil {
		return redis.ErrClosed
	}
	return s.client.Ping(s.ctx).Err()
}

func (s *RedisStorage) dbg(format string, args ...interface{}) {
	if s == nil || !s.debug {
		return
	}
	log.Printf("[identityMng-redis] "+format, args...)
}
