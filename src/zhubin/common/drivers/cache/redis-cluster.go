package cache

import (
	radixCluster "github.com/mediocregopher/radix.v2/cluster"
	radix "github.com/mediocregopher/radix.v2/redis"
	"strings"
	"time"
)

// L2 cache implementation.
// use Radix client because it handles cluster structure transparently, though it is not as type-safe as go-redis client.
type RedisCluster struct {
	cluster *radixCluster.Cluster
}

func NewRedisCluster(address string, password string) *RedisCluster {
	env_password = password
	opts := radixCluster.Opts{
		Addr:     address,
		Timeout:  time.Second * 5,
		PoolSize: 10,
		Dialer:   dialFunc,
	}
	cs, err := radixCluster.NewWithOpts(opts)
	if err != nil {
		panic(err)
	}
	return &RedisCluster{cluster: cs}
}

var env_password string

// helper function to dial the redis server: do possible auth
// SELECT is not allowed in cluster mode
func dialFunc(network, address string) (*radix.Client, error) {
	c, err := radix.Dial(network, address)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(env_password) != "" {
		if err := c.Cmd("AUTH", env_password).Err; err != nil {
			c.Close()
			return nil, err
		}
	}
	return c, nil
}

func (s *RedisCluster) Keys(pattern string) ([]string, error) {
	// ref: http://redis.io/commands/KEYS
	return s.cluster.Cmd("KEYS", pattern).List()
}

// Get fetches value by given key from redis.
func (s *RedisCluster) Get(key string) (value string, err error) {
	resp := s.cluster.Cmd("GET", key)
	if resp.Err != nil {
		return "", resp.Err
	}
	// NOTE: bugfix by Jianjun Xie
	// we now report errNil when no data associated with the key instead of nil error with empty string,
	// since we do allow to store empty string as value
	v, err := resp.Str()
	if err == radix.ErrRespNil {
		return "", ErrKeyNotFound
	}
	return v, nil
}

// Get fetches value by given key from redis.
func (s *RedisCluster) GetObject(key string) (value interface{}, err error) {
	resp := s.cluster.Cmd("GET", key)
	if resp.Err != nil {
		return "", resp.Err
	}
	// NOTE: bugfix by Jianjun Xie
	// we now report errNil when no data associated with the key instead of nil error with empty string,
	// since we do allow to store empty string as value
	v, err := resp.Str()
	if err == radix.ErrRespNil {
		return "", ErrKeyNotFound
	}
	return v, nil
}

// Stores the key and corresponding value in Redis.
func (s *RedisCluster) Save(key string, value string, exp_time_seconds int) error {
	if exp_time_seconds <= 0 {
		return s.cluster.Cmd("SET", key, value).Err
	}
	return s.cluster.Cmd("SETEX", key, exp_time_seconds, value).Err
}

// Stores the key and corresponding value in Redis.
func (s *RedisCluster) SaveObject(key string, value interface{}, exp_time_seconds int) error {
	if exp_time_seconds <= 0 {
		return s.cluster.Cmd("SET", key, value).Err
	}
	return s.cluster.Cmd("SETEX", key, exp_time_seconds, value).Err
}

func (s *RedisCluster) SaveIfNotExists(key string, value string, exp_time_seconds int) error {
	// use new command since 2.6.12
	// see: https://redis.io/commands/set
	return s.cluster.Cmd("set", key, value, "ex", exp_time_seconds, "nx").Err
}

func (s *RedisCluster) SaveObjectIfNotExists(key string, value interface{}, exp_time_seconds int) error {
	// use new command since 2.6.12
	// see: https://redis.io/commands/set
	return s.cluster.Cmd("set", key, value, "ex", exp_time_seconds, "nx").Err
}

// check if key exists
func (s *RedisCluster) Exists(key string) (bool, error) {
	resp := s.cluster.Cmd("EXISTS", key)
	v, err := resp.Int()
	return v != 0, err
}

// Set the expire
func (s *RedisCluster) SetExpire(key string, exp_time_seconds int) error {
	return s.cluster.Cmd("EXPIRE", key, exp_time_seconds).Err
}

// Increments the number stored at key by one.
func (s *RedisCluster) Incr(key string) (int64, error) {
	return s.cluster.Cmd("INCR", key).Int64()
}

func (s *RedisCluster) IncrBy(key string, incValue int64) (int64, error) {
	return s.cluster.Cmd("INCRBY", key, incValue).Int64()
}

func (s *RedisCluster) Delete(key ...string) error {
	return s.cluster.Cmd("DEL", key).Err
}

func (s *RedisCluster) Close() error {
	s.cluster.Close()
	return nil
}