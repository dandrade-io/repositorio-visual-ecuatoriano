package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"uuid"
)

type Asset struct {
	ID         uuid.UUID
	Path       string
	Kind       string
	UploadedBy uuid.UUID
}

func createAssetQuery(a Asset) (string, []any) {
	return "INSERT INTO asset (id, path, kind, uploaded_by) VALUES($1, $2, $3, $4)", []any{a.ID, a.Path, a.Kind, a.UploadedBy}
}

func InsertAsset(ctx context.Context, a Asset) error {
	q, args := createAssetQuery(a)
	_, err := db.ExecContext(ctx, q, args...)

	return err
}

func InsertAssetTx(ctx context.Context, tx *sql.Tx, a Asset) error {
	q, args := createAssetQuery(a)
	_, err := tx.ExecContext(ctx, q, args...)

	return err
}

// ClaimAssetTx sets asset's pending field to True
func ClaimAssetTx(ctx context.Context, tx *sql.Tx, assetID uuid.UUID, userID uuid.UUID) error {
	_, err := tx.ExecContext(
		ctx,
		`UPDATE asset SET pending = false
		WHERE id = $1 AND uploaded_by = $2 AND pending`,
		assetID, userID,
	)

	return err
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
