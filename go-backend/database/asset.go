package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"uuid"
)

type Asset struct {
	ID   uuid.UUID
	Path string
	Kind string
}

func CreateAsset(a Asset) error {
	_, err := db.Exec("INSERT INTO asset (id, path, kind) VALUES($1, $2, $3)", a.ID, a.Path, a.Kind)

	if err != nil {
		return err
	}

	return nil
}

func GetAssetByID(ctx context.Context, id uuid.UUID) (Asset, error) {
	const query = `SELECT id, path, kind FROM asset WHERE id = $1`

	var a Asset
	if err := db.QueryRowContext(ctx, query, id).Scan(&a.ID, &a.Path, &a.Kind); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Asset{}, ErrNotFound
		}
		return Asset{}, fmt.Errorf("get asset by id: %w", err)
	}

	return a, nil
}
