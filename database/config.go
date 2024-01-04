package database

import (
	"encoding/json"
	"fmt"
)

type Config struct {
	Driver   string
	Host     string
	Port     json.Number
	Database string
	Username string
	Password string
	SSLMode  string
}

func (c Config) Dsn() string {
	switch c.Driver {
	case "mysql":
		return c.mysqlDsn()
	case "postgres":
		return c.postgresDsn()
	case "sqlite", "sqlite3":
		return c.Database
	default:
		return ""
	}
}

func (c *Config) mysqlDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4",
		c.Username, c.Password, c.Host, c.Port, c.Database)
}

func (c *Config) postgresDsn() string {
	if c.SSLMode == "" {
		c.SSLMode = "disable"
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Password, c.Database, c.SSLMode)
}
