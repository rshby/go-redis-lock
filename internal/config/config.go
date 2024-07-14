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

// MySqlHost is function to get db mysql host from env
func MySqlHost() string {
	mode := Mode()
	switch mode {
	case ModeLocal:
		if host := GetEnv("DB_HOST"); host != "" {
			return host
		}

		return DefaultMySqlHost
	case ModeDev:
		if host := GetEnv("DB_LOCAL_HOST"); host != "" {
			return host
		}

		return DefaultMySqlHost
	default:
		return DefaultMySqlHost
	}
}

// MySqlPort is function to get mysql port from env
func MySqlPort() int {
	if port := GetEnv("DB_PORT"); port != "" {
		mysqlPort, err := strconv.Atoi(port)
		if err != nil {
			return DefaultMySqlPort
		}

		return mysqlPort
	}

	return DefaultMySqlPort
}

// MySqlUser is function to get mysql user from env
func MySqlUser() string {
	if user := GetEnv("DB_USER"); user != "" {
		return user
	}

	return DefaultMySqlUser
}

// MySqlPassword is function to get mysql password from env
func MySqlPassword() string {
	if password := GetEnv("DB_PASSWORD"); password != "" {
		return password
	}

	return DefaultMySqlPassword
}

// MySqlDbName is function to get db name from env
func MySqlDbName() string {
	if name := GetEnv("DB_NAME"); name != "" {
		return name
	}

	return DefaultMysqlDbName
}

// MySqlDSN is function to connect with database mysql
func MySqlDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=UTC", MySqlUser(), MySqlPassword(), MySqlHost(), MySqlPort(), MySqlDbName())
}

// EnableMigrationDbMysql is function to get enable migration from env
func EnableMigrationDbMysql() bool {
	if enable := GetEnv("ENABLE_MIGRATION_DB"); enable != "" {
		parseBool, err := strconv.ParseBool(enable)
		if err != nil {
			return DefaultEnableMigration
		}

		return parseBool
	}

	return DefaultEnableMigration
}

// MySqlMaxIdleConns is function to get mysql max idle conns from env
func MySqlMaxIdleConns() int {
	if conn := GetEnv("DB_MAX_IDLE_CONNS"); conn != "" {
		mysqlMaxIdleConns, err := strconv.Atoi(conn)
		if err != nil {
			return DefaultMySqlIdleConns
		}

		return mysqlMaxIdleConns
	}

	return DefaultMySqlIdleConns
}

// MySqlMaxOpenConns is function to get max open conns from env
func MySqlMaxOpenConns() int {
	if conns := GetEnv("DB_MAX_OPEN_CONNS"); conns != "" {
		maxOpenConns, err := strconv.Atoi(conns)
		if err != nil {
			return DefaultMysqlMaxOpenConns
		}

		return maxOpenConns
	}

	return DefaultMysqlMaxOpenConns
}

// MysqlConnMaxIdletime is function to get conn max idletime from env
func MysqlConnMaxIdletime() time.Duration {
	if tm := GetEnv("DB_CONN_MAX_IDLETIME"); tm != "" {
		duration, err := time.ParseDuration(tm)
		if err != nil {
			return DefaultMySqlConnMaxIdletime
		}

		return duration
	}

	return DefaultMySqlConnMaxIdletime
}

// MySqlConnMaxLifetime is function to get conn max lifetime from env
func MySqlConnMaxLifetime() time.Duration {
	if lifetime := GetEnv("DB_CONN_MAX_LIFETIME"); lifetime != "" {
		duration, err := time.ParseDuration(lifetime)
		if err != nil {
			return DefaultMySqlConnMaxLifetime
		}

		return duration
	}

	return DefaultMySqlConnMaxLifetime
}
