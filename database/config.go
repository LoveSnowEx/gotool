package database

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Driver   string      `mapstructure:"driver"`
	Host     string      `mapstructure:"host"`
	Port     json.Number `mapstructure:"port"`
	Database string      `mapstructure:"database"`
	Username string      `mapstructure:"username"`
	Password string      `mapstructure:"password"`
	SSLMode  string      `mapstructure:"sslmode"`
}

func ReadViper(v *viper.Viper) (c *Config, err error) {
	c = new(Config)
	err = v.Unmarshal(c)
	if err != nil {
		panic(err)
	}
	return
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
