package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rshby/go-redis-lock/database"
	"github.com/rshby/go-redis-lock/internal/cache"
	"github.com/rshby/go-redis-lock/internal/config"
	"github.com/rshby/go-redis-lock/internal/logger"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
)

func init() {
	logger.SetupLogger()
}

func main() {
	// initialize redis
	database.InitializeRedisConn(config.RedisDSN(), nil)

	// initialize cacheManager
	cacheManager := cache.NewCacheManager(database.RedisConnPool)
	cacheManager.Set("ok", 1)

	// initialize app
	app := gin.Default()

	// server
	server := http.Server{
		Addr:    ":4000",
		Handler: app,
	}

	var (
		chanSignal      = make(chan os.Signal, 1)
		chanServerError = make(chan error, 1)
		chanExit        = make(chan bool)
	)

	signal.Notify(chanSignal, os.Interrupt)

	//
	go func() {
		for {
			select {
			case <-chanSignal:
				// TODO : gracefull shutdown
				logrus.Info("receive interupt signal")
				chanExit <- true
			case err := <-chanServerError:
				logrus.Error(err)
				// TODO : gracefull shutdown
				chanExit <- true
			case <-chanExit:
				return
			}
		}
	}()

	// run server
	go func() {
		logrus.Infof("app run in port 4000")
		if err := server.ListenAndServe(); err != nil {
			chanServerError <- err
			return
		}
	}()

	<-chanExit
	close(chanExit)
	logrus.Info("server exit")
}
