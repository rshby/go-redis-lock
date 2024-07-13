package main

import (
	"encoding/json"
	"fmt"
	"github.com/rshby/go-redis-lock/database"
	"github.com/rshby/go-redis-lock/internal/cache"
	"github.com/rshby/go-redis-lock/internal/config"
	"github.com/rshby/go-redis-lock/internal/logger"
	"github.com/sirupsen/logrus"
	"reflect"
	"time"
)

func init() {
	logger.SetupLogger()
}

type User struct {
	ID   int
	Name string
}

func main() {
	// initialize redis
	database.InitializeRedisConn(config.RedisDSN(), nil)

	// initialize cacheManager
	cacheManager := cache.NewCacheManager(database.RedisConnPool)

	user1 := User{1, "Febian Diska Haryuningtyas"}
	user1Json, _ := json.Marshal(&user1)
	if err := cacheManager.Set("oke3", string(user1Json)); err != nil {
		logrus.Error(err)
	}

	key, err := cache.GetByKey[*User](cacheManager, "oke3")
	if err != nil {
		logrus.Error(err)
	}

	fmt.Println(key, reflect.TypeOf(key))

	time.Sleep(1 * time.Hour)
}
