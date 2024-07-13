package cache

import (
	"encoding/json"
	"github.com/go-redsync/redsync/v4"
	redigosync "github.com/go-redsync/redsync/v4/redis/redigo"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/rshby/go-redis-lock/internal/cache/interfaces"
	"github.com/rshby/go-redis-lock/internal/config"
	"reflect"
	"time"
)

type cacheManager struct {
	environment    string
	prefixCacheKey string

	connPool    *redigo.Pool
	nilTTL      time.Duration
	defaultTTL  time.Duration
	waitTime    time.Duration
	enableCache bool

	lockConnPool *redigo.Pool
	lockDuration time.Duration
	lockTries    int
}

// NewCacheManager is function to create new instance cacheManager
func NewCacheManager(redisConnPool *redigo.Pool) interfaces.CacheManager {
	return &cacheManager{
		environment:    config.Mode(),
		prefixCacheKey: config.RedisPrefixKey(),
		connPool:       redisConnPool,
		defaultTTL:     config.RedisTTL(),
		enableCache:    config.EnableCache(),
		lockConnPool:   redisConnPool,
	}
}

// Get is method to get data from cache
func (c *cacheManager) Get(key string) ([]byte, error) {
	if !c.enableCache {
		return nil, nil
	}

	// create client
	client := c.connPool.Get()
	defer client.Close()

	key = c.prefixCacheKey + key
	reply, err := client.Do("GET", key)
	if err != nil {
		return nil, err
	}
	if reply == nil {
		return nil, nil
	}

	return reply.([]byte), nil
}

// Set is method to set data to cache
func (c *cacheManager) Set(key string, value any) error {
	if !c.enableCache {
		return nil
	}

	// TODO : Lock

	// create client
	client := c.connPool.Get()
	defer client.Close()

	key = c.prefixCacheKey + key
	_, err := client.Do("SETEX", key, int64(c.defaultTTL.Seconds()), value)
	if err != nil {
		return err
	}

	return nil
}

// AcquireLock is function to acquire lock by key
func (c *cacheManager) AcquireLock(key string) (*redsync.Mutex, error) {
	key = c.prefixCacheKey + "lock:" + key

	// create new pool
	pool := redigosync.NewPool(c.lockConnPool)

	// create mutex
	mutex := redsync.New(pool).NewMutex(
		key,
		redsync.WithExpiry(30*time.Second),
		redsync.WithTries(3),
	)

	// lock the mutex
	return mutex, mutex.Lock()
}

// SafeUnlock is function to unlock mutex
func (c *cacheManager) SafeUnlock(mutex *redsync.Mutex) {
	if mutex != nil {
		_, _ = mutex.Unlock()
	}
}

// GetByKey is function to get value by key
func GetByKey[T any | string](c interfaces.CacheManager, key string) (T, error) {
	var result T

	response, err := c.Get(key)
	if err != nil {
		return result, err
	}

	if response == nil {
		return result, nil
	}

	switch reflect.TypeOf(result).Kind() {
	case reflect.String:
		responseString := string(response)
		responseBytes, err := json.Marshal(&responseString)
		if err != nil {
			return result, err
		}

		if err = json.Unmarshal(responseBytes, &result); err != nil {
			return result, err
		}
		return result, nil
	default:
		if err = json.Unmarshal(response, &result); err != nil {
			return result, err
		}

		return result, nil
	}
}
