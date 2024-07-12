package main

import (
	"github.com/rshby/go-redis-lock/database"
	"github.com/rshby/go-redis-lock/internal/config"
	"github.com/rshby/go-redis-lock/internal/logger"
	"time"
)

func init() {
	logger.SetupLogger()
}

func main() {
	database.InitializeRedisConn(config.RedisDSN(), nil)

	time.Sleep(1 * time.Hour)
}
