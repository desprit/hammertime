package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"desprit/hammertime/src/config"
)

func Open(cfg *config.Config) (*sql.DB, error) {
	d, err := sql.Open("sqlite3", cfg.DbUri())
	if err != nil {
		return nil, err
	}
	if err := d.Ping(); err != nil {
		return nil, err
	}
	d.SetMaxOpenConns(1)
	d.SetMaxIdleConns(1)
	return d, nil
}

func CreateTables(cfg *config.Config, d *sql.DB) error {
	tables := []string{"schedule", "subscription"}
	for _, tableName := range tables {
		query, err := Asset(fmt.Sprintf("src/db/%s/schema.sql", tableName))
		if err != nil {
			return err
		}
		_, err = d.Exec(string(query))
		if err != nil {
			return err
		}
	}

	return nil
}

func InitDB(d *sql.DB, cfg *config.Config) error {
	err := CreateTables(cfg, d)
	if err != nil {
		return err
	}

	return nil
}
