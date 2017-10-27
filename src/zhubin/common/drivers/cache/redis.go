package cache

import (
	"github.com/go-redis/redis"
	"time"
)

// L2 cache implementation.
type Redis struct {
	client *redis.Client
}

// create a redis client.
func NewRedis(address string, password string, dbnum int) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       dbnum,
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic("Failed to connect redis.")
	}
	rs := &Redis{
		client: client,
	}
	return rs
}

func (s *Redis) Keys(pattern string) ([]string, error) {
	return s.client.Keys(pattern).Result()
}

// Get fetches value by given key from redis.
func (s *Redis) Get(key string) (value string, err error) {
	val, err := s.client.Get(key).Result()
	// NOTE: bugfix by Jianjun Xie
	// we now report errNil when no data associated with the key instead of nil error with empty string,
	// since we do allow to store empty string as value
	if err == redis.Nil {
		return val, ErrKeyNotFound
	}
	return val, err
}

func (s *Redis) Save(key string, value string, exp_time_seconds int) error {
	return s.client.Set(key, value, time.Second*time.Duration(exp_time_seconds)).Err()
}

func (s *Redis) SaveObject(key string, value interface{}, exp_time_seconds int) error {
	return s.client.Set(key, value, time.Second*time.Duration(exp_time_seconds)).Err()
}

func (s *Redis) SaveIfNotExists(key string, value string, exp_time_seconds int) error {
	return s.client.SetNX(key, value, time.Second*time.Duration(exp_time_seconds)).Err()
}

func (s *Redis) SaveObjectIfNotExists(key string, value interface{}, exp_time_seconds int) error {
	return s.client.SetNX(key, value, time.Second*time.Duration(exp_time_seconds)).Err()
}

func (s *Redis) Exists(key string) (bool, error) {
	b, err := s.client.Exists(key).Result()
	return b != 0, err
}

func (s *Redis) SetExpire(key string, exp_time_seconds int) error {
	return s.client.Expire(key, time.Second*time.Duration(exp_time_seconds)).Err()
}

func (s *Redis) Incr(key string) (int64, error) {
	return s.client.Incr(key).Result()
}

func (s *Redis) IncrBy(key string, incValue int64) (int64, error) {
	return s.client.IncrBy(key, incValue).Result()
}

func (s *Redis) Delete(key ...string) error {
	return s.client.Del(key...).Err()
}

func (s *Redis) Close() error {
	return s.client.Close()
}
