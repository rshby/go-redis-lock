package database

import (
	"github.com/go-redis/redis/v7"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/rshby/go-redis-lock/internal/config"
	"github.com/sirupsen/logrus"
	"time"
)

type RedisConnectionPoolOptions struct {
	DialTimeout     time.Duration
	ReadOnly        bool
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleCount       int
	PoolSize        int
	IdleTimeout     time.Duration
	MaxConnLifetime time.Duration
}

var (
	RedisConnPool                     *redigo.Pool
	StopTicker                        = make(chan bool)
	defaultRedisConnectionPoolOptions = RedisConnectionPoolOptions{
		DialTimeout:     config.DefaultRedisTimeout,
		ReadOnly:        config.DefaultRedisReadOnly,
		ReadTimeout:     config.DefaultRedisReadTimeout,
		WriteTimeout:    config.DefaultRedisWriteTimeout,
		IdleCount:       config.DefaultRedisIdleCount,
		PoolSize:        config.DefaultRedisPoolSize,
		IdleTimeout:     config.DefaultRedisIdleTimeout,
		MaxConnLifetime: config.DefaultRedisMaxConnLifetime,
	}
)

// InitializeRedisConn is function to connect redis
func InitializeRedisConn(url string, opt *RedisConnectionPoolOptions) {
	logger := logrus.WithFields(logrus.Fields{
		"url": url,
	})

	// validate redis url
	if !IsValidURL(url) {
		logrus.WithField("urlRedis", url).Error("cant connect redis")
	}

	// apply redis connection pool options
	opt = ApplyRedisConnectionPoolOptions(opt)

	// create redis client
	RedisConnPool = &redigo.Pool{
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.DialURL(url)
			if err != nil {
				logger.Error(err)
				return nil, err
			}

			return c, nil
		},
		TestOnBorrow: func(c redigo.Conn, lastUsed time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				logger.Error(err)
			}

			return err
		},
		MaxIdle:         opt.IdleCount,
		MaxActive:       opt.PoolSize,
		IdleTimeout:     opt.IdleTimeout,
		Wait:            true,
		MaxConnLifetime: opt.MaxConnLifetime,
	}

	go CheckRedisContinously(RedisConnPool, time.NewTicker(5*time.Second), url, opt)

	logger.Info("success create connection redis ", url)
}

// IsValidURL is function to check url is valid or not
func IsValidURL(url string) bool {
	_, err := redis.ParseURL(url)
	if err != nil {
		logrus.WithField("urlRedis", url).Error(err)
		return false
	}

	return true
}

// ApplyRedisConnectionPoolOptions is method to apply redis connection pool options
func ApplyRedisConnectionPoolOptions(opt *RedisConnectionPoolOptions) *RedisConnectionPoolOptions {
	if opt == nil {
		return &defaultRedisConnectionPoolOptions
	}

	return opt
}

// CheckRedisContinously is function to
func CheckRedisContinously(redisConnPool *redigo.Pool, ticker *time.Ticker, url string, opt *RedisConnectionPoolOptions) {
	for {
		select {
		case <-StopTicker:
			ticker.Stop()
			return
		case <-ticker.C:
			client := redisConnPool.Get()
			if _, err := client.Do("PING"); err != nil {
				// reconnect redis
				ReconnectRedis(url, opt)
			}

			if err := client.Close(); err != nil {
				logrus.Error(err)
			}
		}
	}
}

// ReconnectRedis is function to reconnect redis
func ReconnectRedis(url string, opt *RedisConnectionPoolOptions) {
	logger := logrus.WithFields(logrus.Fields{
		"url": url,
	})

	// apply redis connection pool options
	opt = ApplyRedisConnectionPoolOptions(opt)

	// create redis connection pool
	redisConnPool := redigo.Pool{
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.DialURL(url)
			if err != nil {
				logger.Error(err)
				return nil, err
			}

			return c, nil
		},
		TestOnBorrow: func(c redigo.Conn, lastUsed time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				logger.Error(err)
			}

			return err
		},
		MaxIdle:         opt.IdleCount,
		MaxActive:       opt.PoolSize,
		IdleTimeout:     opt.IdleTimeout,
		Wait:            true,
		MaxConnLifetime: opt.MaxConnLifetime,
	}

	// create client
	c := redisConnPool.Get()
	if _, err := c.Do("PING"); err == nil {
		RedisConnPool = &redisConnPool
	}

	// close client
	_ = c.Close()
}
