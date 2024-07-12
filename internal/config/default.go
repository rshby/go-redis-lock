package config

import "time"

const (
	// app
	DefaultMode = "local"
	ModeLocal   = "local"
	ModeDev     = "dev"

	// redis
	DefaultRedisHost            = "localhost"
	DefaultRedisPort            = 6379
	DefaultRedisDBNumber        = 1
	DefaultRedisTimeout         = 5 * time.Second
	DefaultRedisReadOnly        = false
	DefaultRedisReadTimeout     = 2 * time.Second
	DefaultRedisWriteTimeout    = 2 * time.Second
	DefaultRedisIdleCount       = 20
	DefaultRedisPoolSize        = 100
	DefaultRedisIdleTimeout     = 1 * time.Minute
	DefaultRedisMaxConnLifetime = 1 * time.Hour
)
