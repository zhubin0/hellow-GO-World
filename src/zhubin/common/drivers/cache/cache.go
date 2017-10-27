package cache

import (
	"errors"
)

var ErrKeyNotFound = errors.New("key not found")

type L2Cache interface {
	// Returns all keys matching pattern.
	Keys(pattern string) ([]string, error)

	// Get fetches value by given key from cache.
	// Should report ErrKeyNotFound when no associated value found
	Get(key string) (value string, err error)
	//GetObject(key string) (value interface{}, err errors)

	// Save stores the key and corresponding value in Redis.
	// if exp_time_seconds <= 0, the key never expire
	Save(key string, value string, exp_time_seconds int) error

	// Save stores the key， value is inerface{}
	SaveObject(key string, value interface{}, exp_time_seconds int) error

	// Set key to hold string value if key does not exist.
	SaveIfNotExists(key string, value string, exp_time_seconds int) error
	SaveObjectIfNotExists(key string, value interface{}, exp_time_seconds int) error

	// Check if the given key exists.
	Exists(key string) (bool, error)

	// Set the expire
	SetExpire(key string, exp_time_seconds int) error

	// Increments the number stored at key by one, return the value after increment.
	Incr(key string) (int64, error)

	// Increments the number stored at key by given incValue, return the value after increment.
	IncrBy(key string, incValue int64) (int64, error)

	// Delete removes one or more keys.
	Delete(key ...string) error

	// close the cache client
	Close() error
}