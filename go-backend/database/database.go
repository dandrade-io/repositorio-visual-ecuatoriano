package database

import (
	"database/sql"
	"errors"
	"fmt"
	"repository/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB
var ErrNotFound = errors.New("record not found")

func Open() error {
	var err error
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	db, err = sql.Open("pgx", dbURL)

	return err
}

func Close() error {
	return db.Close()
}
