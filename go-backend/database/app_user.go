package database

import (
	"context"
	"uuid"
)

func InsertUser(ctx context.Context, id uuid.UUID) error {
	_, err := db.ExecContext(
		ctx,
		`INSERT INTO app_user(id) VALUES ($1) ON CONFLICT(id) DO NOTHING`,
		id,
	)

	return err
}
