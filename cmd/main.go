package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rshby/go-redis-lock/database"
	"github.com/rshby/go-redis-lock/http/router"
	"github.com/rshby/go-redis-lock/internal/cache"
	"github.com/rshby/go-redis-lock/internal/config"
	"github.com/rshby/go-redis-lock/internal/logger"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func init() {
	logger.SetupLogger()
}

func main() {
	// initialize db mysql
	database.InitializeMysql()

	// initialize redis
	database.InitializeRedisConn(config.RedisDSN(), nil)

	// initialize cacheManager
	cacheManager := cache.NewCacheManager(database.RedisConnPool)

	// initialize app
	app := gin.Default()

	// router
	appRouter := router.NewAppRouter(&app.RouterGroup, cacheManager)
	appRouter.InitEndpoint()

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
				logrus.Info("receive interupt signal")
				GracefullShutdown(&server)
				chanExit <- true
				return
			case err := <-chanServerError:
				logrus.Error(err)
				GracefullShutdown(&server)
				chanExit <- true
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
	logrus.Info("server exit❌")
}

// GracefullShutdown is function to shut down server
func GracefullShutdown(server *http.Server) {
	// stop and shutdown server
	if server != nil {
		server.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logrus.Fatal("force close server🔴")
		}

		logrus.Info("success shutdown server")
	}

	// stop database connection
	database.StopTicker <- true
	time.Sleep(500 * time.Millisecond)
	database.CloseMySqlConnection(database.DatabaseMySQL)
	database.CloseConnection(database.RedisConnPool)
}
