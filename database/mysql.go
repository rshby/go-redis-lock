package database

import (
	"github.com/rshby/go-redis-lock/database/migration"
	"github.com/rshby/go-redis-lock/internal/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	dbMysql.SetMaxIdleConns(config.MySqlMaxIdleConns())
	dbMysql.SetMaxOpenConns(config.MySqlMaxOpenConns())
	dbMysql.SetConnMaxIdleTime(config.MysqlConnMaxIdletime())
	dbMysql.SetConnMaxLifetime(config.MySqlConnMaxLifetime())

	return db
}

// CloseMySqlConnection is function to close db mysql connection
func CloseMySqlConnection(db *gorm.DB) {
	if db == nil {
		return
	}

	mySqlDB, err := db.DB()
	if err != nil {
		logrus.Error(err)
		return
	}

	if err := mySqlDB.Close(); err != nil {
		logrus.Infof("failed to close db mysql connection : %v", err)
		return
	}

	logrus.Info("succes close mysql connection")
}
