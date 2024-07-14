package config

import "time"

const (
	// app
	DefaultAppName = "go-redis-lock"
	DefaultMode    = "local"
	ModeLocal      = "local"
	ModeDev        = "dev"

	// redis
	DefaultEnableCache          = false
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
	DefaultRedisTTL             = 15 * time.Minute

	// mysql
	DefaultMySqlHost            = "localhost"
	DefaultMySqlPort            = 3306
	DefaultMySqlUser            = "root"
	DefaultMySqlPassword        = "root"
	DefaultMysqlDbName          = "go_redis_lock_db"
	DefaultMySqlIdleConns       = 30
	DefaultMysqlMaxOpenConns    = 100
	DefaultMySqlConnMaxIdletime = 30 * time.Minute
	DefaultMySqlConnMaxLifetime = 1 * time.Hour
	DefaultEnableMigration      = false
)
