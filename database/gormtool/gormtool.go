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

func Open(cfg *database.Config, opts ...gorm.Option) (db *gorm.DB, err error) {
	switch cfg.Driver {
	case "mysql":
		return mysqlOpen(cfg, opts...)
	case "postgres":
		return postgresOpen(cfg, opts...)
	case "sqlite", "sqlite3":
		return sqliteOpen(cfg, opts...)
	default:
		return nil, errors.ErrUnknownDatabaseDriver
	}
}

func mysqlOpen(cfg *database.Config, opts ...gorm.Option) (db *gorm.DB, err error) {
	return gorm.Open(mysql.Open(cfg.Dsn()), opts...)
}

func postgresOpen(cfg *database.Config, opts ...gorm.Option) (db *gorm.DB, err error) {
	return gorm.Open(postgres.Open(cfg.Dsn()), opts...)
}

func sqliteOpen(cfg *database.Config, opts ...gorm.Option) (db *gorm.DB, err error) {
	return gorm.Open(sqlite.Open(cfg.Dsn()), opts...)
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
