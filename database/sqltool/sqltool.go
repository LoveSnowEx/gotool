package sqltool

import (
	"database/sql"

	"github.com/LoveSnowEx/gotool/database"
	sq "github.com/Masterminds/squirrel"
)

func Open(cfg database.Config) (db *sql.DB, err error) {
	return sql.Open(cfg.Driver, cfg.Dsn())
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

func Tables(db *sql.DB) (tables []string, err error) {
	sql, args, err := sq.
		Select("TABLE_NAME").
		From("INFORMATION_SCHEMA.TABLES").
		Where("TABLE_SCHEMA = DATABASE()").
		ToSql()
	if err != nil {
		return
	}

	rows, err := db.Query(sql, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var table string
		if err = rows.Scan(&table); err != nil {
			return
		}
		tables = append(tables, table)
	}
	return
}
