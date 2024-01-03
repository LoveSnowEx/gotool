package sqltool

import (
	"database/sql"

	"github.com/LoveSnowEx/gotool/database"
	sq "github.com/Masterminds/squirrel"
)

func Open(conf *database.DatabaseConfig) (db *sql.DB, err error) {
	return sql.Open(conf.Driver, conf.Dsn())
}

func PK(db *sql.DB, table string) (pk string, err error) {
	sql, args, err := sq.
		Select("COLUMN_NAME").
		From("INFORMATION_SCHEMA.KEY_COLUMN_USAGE").
		Where("CONSTRAINT_NAME = ? AND TABLE_NAME = ? AND TABLE_SCHEMA = DATABASE()", "PRIMARY", table).
		ToSql()
	if err != nil {
		return
	}

	if err = db.QueryRow(sql, args...).Scan(&pk); err != nil {
		return
	}
	return
}
