package gormtool

import (
	"github.com/LoveSnowEx/gotool/database"
	"github.com/LoveSnowEx/gotool/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Open(conf *database.Config, gormConf *gorm.Config) (db *gorm.DB, err error) {
	switch conf.Driver {
	case "mysql":
		return mysqlOpen(conf, gormConf)
	case "postgres":
		return postgresOpen(conf, gormConf)
	case "sqlite", "sqlite3":
		return sqliteOpen(conf, gormConf)
	default:
		return nil, errors.ErrUnknownDatabaseDriver
	}
}

func mysqlOpen(conf *database.Config, gormConf *gorm.Config) (db *gorm.DB, err error) {
	return gorm.Open(mysql.Open(conf.Dsn()), gormConf)
}

func postgresOpen(conf *database.Config, gormConf *gorm.Config) (db *gorm.DB, err error) {
	return gorm.Open(postgres.Open(conf.Dsn()), gormConf)
}

func sqliteOpen(conf *database.Config, gormConf *gorm.Config) (db *gorm.DB, err error) {
	return gorm.Open(sqlite.Open(conf.Dsn()), gormConf)
}
