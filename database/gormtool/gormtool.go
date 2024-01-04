package gormtool

import (
	"sync"

	"github.com/LoveSnowEx/gotool/database"
	"github.com/LoveSnowEx/gotool/errors"
	"github.com/panjf2000/ants/v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gen"
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

func GenDao(g *gen.Generator, numWorker int, tables ...string) (err error) {
	var wg sync.WaitGroup
	p, err := ants.NewPoolWithFunc(numWorker, func(table interface{}) {
		t := table.(string)
		g.ApplyBasic(g.GenerateModel(t))
		wg.Done()
	})
	if err != nil {
		return
	}
	defer p.Release()

	for _, table := range tables {
		wg.Add(1)
		if err = p.Invoke(table); err != nil {
			return
		}
	}
	wg.Wait()
	g.Execute()
	return
}
