package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"repository/config"

	"github.com/jackc/pgx/v5/pgtype"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB
var ErrNotFound = errors.New("record not found")
var typeMap = pgtype.NewMap()

func Open() error {
	var err error
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	db, err = sql.Open("pgx", dbURL)

	return err
}

func Close() error {
	return db.Close()
}

func WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}
