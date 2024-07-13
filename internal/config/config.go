package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

// GetEnv is function to load env and get value by key
func GetEnv(key string) string {
	logger := logrus.WithField("key", key)

	if err := godotenv.Load(".env"); err != nil {
		logger.Fatal(err)
	}

	return os.Getenv(key)
}

// AppName is method to get app name from env
func AppName() string {
	if name := GetEnv("APP_NAME"); name != "" {
		return name
	}

	return DefaultAppName
}

// Mode is function to get mode from env
func Mode() string {
	mode := GetEnv("MODE")

	if mode != "" {
		return mode
	}

	switch mode {
	case ModeLocal:
		return ModeLocal
	case ModeDev:
		return ModeDev
	default:
		return DefaultMode
	}
}

// EnableCache is method to get enable cache from env
func EnableCache() bool {
	if enablCache := GetEnv("ENABLE_CACHE"); enablCache != "" {
		parseBool, err := strconv.ParseBool(enablCache)
		if err != nil {
			return DefaultEnableCache
		}

		return parseBool
	}

	return DefaultEnableCache
}

// RedisHost is function to get redis host from env
func RedisHost() string {
	switch Mode() {
	case ModeLocal:
		if host := GetEnv("REDIS_LOCAL_HOST"); host != "" {
			return host
		}

		return DefaultRedisHost
	case ModeDev:
		if host := GetEnv("REDIS_HOST"); host != "" {
			return host
		}

		return DefaultRedisHost
	default:
		return DefaultRedisHost
	}
}

// RedisPort is function to get redis port from env
func RedisPort() int {
	if port := GetEnv("REDIS_PORT"); port != "" {
		redisPort, err := strconv.Atoi(port)
		if err != nil {
			return DefaultRedisPort
		}

		return redisPort
	}

	return DefaultRedisPort
}

// RedisDBNumber is function to get redis db number from env
func RedisDBNumber() int {
	if num := GetEnv("REDIS_DB_NUMBER"); num != "" {
		redisDBNumber, err := strconv.Atoi(num)
		if err != nil {
			return DefaultRedisDBNumber
		}

		return redisDBNumber
	}

	return DefaultRedisDBNumber
}

// RedisTTL is function to get redis TTL from env
func RedisTTL() time.Duration {
	if ttl := GetEnv("REDIS_TTL"); ttl != "" {
		duration, err := time.ParseDuration(ttl)
		if err != nil {
			return DefaultRedisTTL
		}

		return duration
	}

	return DefaultRedisTTL
}

// RedisPrefixKey is function to get redis prefix key
func RedisPrefixKey() string {
	appName := AppName()
	mode := Mode()
	return fmt.Sprintf("%s_mode:%s_", appName, mode)
}

// RedisDSN is function to get redis dsn
func RedisDSN() string {
	return fmt.Sprintf("redis://%s:%d/%d", RedisHost(), RedisPort(), RedisDBNumber())
}
