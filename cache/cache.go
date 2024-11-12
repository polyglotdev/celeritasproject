package cache

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

type Cacher interface {
	Has(string) (bool, error)
	Get(string) (any, error)
	Set(string, any, ...int) error
	Forget(string) error
	EmptyByMatch(string) error
	Empty() error
}

// RedisCache is a cache implementation using Redis.
type RedisCache struct {
	Pool   *redis.Pool
	Prefix string
}

// Entry is a map of key-value pairs.
type Entry map[string]any

// Has checks if a key exists in the cache.
// It passes in a key and returns a boolean and an error.
func (c *RedisCache) Has(key string) (bool, error) {
	key = fmt.Sprintf("%s%s", c.Prefix, key)
	conn := c.Pool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", c.Prefix+key))
	if err != nil {
		return false, err
	}

	return exists, nil
}
