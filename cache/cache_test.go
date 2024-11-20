package cache

import (
	"strings"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

func setupTestRedis(t *testing.T) *RedisCache {
	pool := &redis.Pool{
		MaxIdle:     50,
		MaxActive:   10000,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}

	// Test the connection
	conn := pool.Get()
	defer conn.Close()

	if _, err := conn.Do("PING"); err != nil {
		t.Skip("Redis server is not running - skipping tests")
	}

	return &RedisCache{
		Conn:   pool,
		Prefix: "test:",
	}
}

func cleanupTestRedis(t *testing.T, cache *RedisCache) {
	if err := cache.Empty(); err != nil {
		t.Errorf("Failed to cleanup test cache: %v", err)
	}
}

func TestRedisCache_Has(t *testing.T) {
	cache := setupTestRedis(t)
	defer cleanupTestRedis(t, cache)

	tests := []struct {
		name    string
		key     string
		value   interface{}
		want    bool
		wantErr bool
	}{
		{
			name:    "existing key",
			key:     "test_has_key",
			value:   "test value",
			want:    true,
			wantErr: false,
		},
		{
			name:    "non-existing key",
			key:     "nonexistent",
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != nil {
				if err := cache.Set(tt.key, tt.value); err != nil {
					t.Fatalf("Failed to set up test: %v", err)
				}
			}

			got, err := cache.Has(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisCache.Has() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RedisCache.Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisCache_Get(t *testing.T) {
	cache := setupTestRedis(t)
	defer cleanupTestRedis(t, cache)

	tests := []struct {
		name    string
		key     string
		value   interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name:    "string value",
			key:     "test_get_string",
			value:   "test value",
			want:    "test value",
			wantErr: false,
		},
		{
			name:    "integer value",
			key:     "test_get_int",
			value:   42,
			want:    42,
			wantErr: false,
		},
		{
			name:    "non-existing key",
			key:     "nonexistent",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != nil {
				if err := cache.Set(tt.key, tt.value); err != nil {
					t.Fatalf("Failed to set up test: %v", err)
				}
			}

			got, err := cache.Get(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisCache.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("RedisCache.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisCache_Set(t *testing.T) {
	cache := setupTestRedis(t)
	defer cleanupTestRedis(t, cache)

	tests := []struct {
		name    string
		key     string
		value   interface{}
		expires []int
		wantErr bool
	}{
		{
			name:    "simple string",
			key:     "test_set_string",
			value:   "test value",
			wantErr: false,
		},
		{
			name:    "with expiration",
			key:     "test_set_expire",
			value:   "expiring value",
			expires: []int{1},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cache.Set(tt.key, tt.value, tt.expires...)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisCache.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify the value was set
			got, err := cache.Get(tt.key)
			if err != nil {
				t.Errorf("Failed to verify set value: %v", err)
				return
			}
			if got != tt.value {
				t.Errorf("Set value = %v, want %v", got, tt.value)
			}

			if len(tt.expires) > 0 {
				// Wait for expiration
				time.Sleep(time.Duration(tt.expires[0]+1) * time.Second)
				_, err := cache.Get(tt.key)
				if err == nil {
					t.Error("Expected key to be expired")
				}
			}
		})
	}
}

func TestRedisCache_Forget(t *testing.T) {
	cache := setupTestRedis(t)
	defer cleanupTestRedis(t, cache)

	// Set up test data
	key := "test_forget"
	if err := cache.Set(key, "test value"); err != nil {
		t.Fatalf("Failed to set up test: %v", err)
	}

	if err := cache.Forget(key); err != nil {
		t.Errorf("RedisCache.Forget() error = %v", err)
	}

	// Verify key was forgotten
	exists, err := cache.Has(key)
	if err != nil {
		t.Errorf("Failed to verify forgotten key: %v", err)
	}
	if exists {
		t.Error("Key should have been forgotten")
	}
}

func TestRedisCache_EmptyByMatch(t *testing.T) {
	cache := setupTestRedis(t)
	defer cleanupTestRedis(t, cache)

	// Set up test data with prefixed keys
	prefix := "test_empty_"
	testData := map[string]string{
		prefix + "1": "value1",
		prefix + "2": "value2",
		"other_key":  "value3",
	}

	for k, v := range testData {
		if err := cache.Set(k, v); err != nil {
			t.Fatalf("Failed to set up test: %v", err)
		}
	}

	// Empty keys matching the pattern
	if err := cache.EmptyByMatch(prefix + "*"); err != nil {
		t.Errorf("RedisCache.EmptyByMatch() error = %v", err)
	}

	// Verify matching keys were removed
	for k := range testData {
		exists, err := cache.Has(k)
		if err != nil {
			t.Errorf("Failed to verify key %s: %v", k, err)
			continue
		}
		if strings.HasPrefix(k, prefix) && exists {
			t.Errorf("Key %s should have been removed", k)
		} else if !strings.HasPrefix(k, prefix) && !exists {
			t.Errorf("Key %s should not have been removed", k)
		}
	}
}

func TestRedisCache_Empty(t *testing.T) {
	cache := setupTestRedis(t)
	defer cleanupTestRedis(t, cache)

	// Set up test data
	testData := map[string]string{
		"test_empty_all_1": "value1",
		"test_empty_all_2": "value2",
	}

	for k, v := range testData {
		if err := cache.Set(k, v); err != nil {
			t.Fatalf("Failed to set up test: %v", err)
		}
	}

	if err := cache.Empty(); err != nil {
		t.Errorf("RedisCache.Empty() error = %v", err)
	}

	// Verify all keys were removed
	for k := range testData {
		exists, err := cache.Has(k)
		if err != nil {
			t.Errorf("Failed to verify key %s: %v", k, err)
			continue
		}
		if exists {
			t.Errorf("Key %s should have been removed", k)
		}
	}
}
