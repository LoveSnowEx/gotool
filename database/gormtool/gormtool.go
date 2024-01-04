package gormtool

import (
	"github.com/LoveSnowEx/gotool/database"
	"github.com/LoveSnowEx/gotool/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Open(conf *database.Config, opts ...gorm.Option) (db *gorm.DB, err error) {
	switch conf.Driver {
	case "mysql":
		return mysqlOpen(conf, opts...)
	case "postgres":
		return postgresOpen(conf, opts...)
	case "sqlite", "sqlite3":
		return sqliteOpen(conf, opts...)
	default:
		return nil, errors.ErrUnknownDatabaseDriver
	}
}

func mysqlOpen(conf *database.Config, opts ...gorm.Option) (db *gorm.DB, err error) {
	return gorm.Open(mysql.Open(conf.Dsn()), opts...)
}

func postgresOpen(conf *database.Config, opts ...gorm.Option) (db *gorm.DB, err error) {
	return gorm.Open(postgres.Open(conf.Dsn()), opts...)
}

func sqliteOpen(conf *database.Config, opts ...gorm.Option) (db *gorm.DB, err error) {
	return gorm.Open(sqlite.Open(conf.Dsn()), opts...)
}
