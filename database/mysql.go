package database

import (
	"github.com/rshby/go-redis-lock/database/migration"
	"github.com/rshby/go-redis-lock/internal/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var (
	DatabaseMySQL *gorm.DB
)

// InitializeMysql is function to connect database
func InitializeMysql() *gorm.DB {
	dsn := config.MySqlDSN()
	DatabaseMySQL = OpenMysqlConnection(dsn)
	logrus.Infof("success connection database mysql %s:%d", config.MySqlHost(), config.MySqlPort())

	// migrations
	if config.EnableMigrationDbMysql() {
		migration.Migration(DatabaseMySQL)
	}

	return DatabaseMySQL
}

// OpenMysqlConnection is function to open connection mysql
func OpenMysqlConnection(url string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(url), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		logrus.Fatal(err)
	}

	dbMysql, err := db.DB()
	if err != nil {
		logrus.Error(err)
	}

	dbMysql.SetMaxIdleConns(20)                  // TODO : ambil dari config
	dbMysql.SetMaxOpenConns(100)                 // TODO : ambil dari config
	dbMysql.SetConnMaxIdleTime(30 * time.Minute) // TODO : ambil dari config
	dbMysql.SetConnMaxLifetime(1 * time.Hour)    // TODO : ambil dari config

	return db
}
