package database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"uuid"
)

type AlbumCover struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Year          *int      `json:"year"`
	Artist        string    `json:"artist"`
	Rhythms       []string  `json:"rhythms"`
	ReprTechnique string    `json:"representationTechnique"`
	Predominance  string    `json:"predominance"`
	Colors        []string  `json:"colors"`
	TitleSyntax   string    `json:"titleSyntax"`
	PublishedBy   uuid.UUID `json:"published_by"`
}

type AlbumCoverDetails struct {
	AlbumCover
	Path   string            `json:"path"`
	Checks []AlbumCoverCheck `json:"checks"`
}

type AlbumCoverSummary struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Year        *int      `json:"year"`
	Artist      string    `json:"artist"`
	PublishedBy uuid.UUID `json:"publishedBy"`
	Path        string    `json:"path"`
}

func QueryAlbumCoverById(ctx context.Context, id uuid.UUID) (AlbumCoverDetails, error) {
	var ac AlbumCoverDetails

	err := db.QueryRowContext(ctx,
		`SELECT DISTINCT
			ac.id, ac.title, ac.year, ac.artist, ac.rhythms, ac.representation_technique, ac.predominance,
		 	ac.colors, ac.title_syntax, ac.published_by, asset.path
		FROM album_cover ac
		INNER JOIN album_cover_asset aca ON aca.album_cover_id = ac.id
		INNER JOIN asset ON aca.asset_id = asset.id
		WHERE ac.id = $1`, id,
	).Scan(&ac.ID, &ac.Title, &ac.Year, &ac.Artist, typeMap.SQLScanner(&ac.Rhythms), &ac.ReprTechnique,
		&ac.Predominance, typeMap.SQLScanner(&ac.Colors), &ac.TitleSyntax, &ac.PublishedBy, &ac.Path,
	)

	if err != nil {
		return AlbumCoverDetails{}, err
	}

	checks, err := QueryAlbumCoverChecksByAlbumCoverID(ctx, id)
	if err != nil {
		return AlbumCoverDetails{}, err
	}
	ac.Checks = checks

	return ac, nil
}

func QueryAlbumCoverSummary(ctx context.Context, limit int) ([]AlbumCoverSummary, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("limit should be greater than 0")
	}

	rows, err := db.QueryContext(
		ctx,
		`SELECT DISTINCT
			ac.id, ac.title, ac.year, ac.artist, ac.published_by, asset.path
		FROM album_cover ac
		INNER JOIN album_cover_asset aca ON aca.album_cover_id = ac.id
		INNER JOIN asset ON aca.asset_id = asset.id
		LIMIT $1
		`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	covers := make([]AlbumCoverSummary, 0, limit)
	for rows.Next() {
		var c AlbumCoverSummary
		if err := rows.Scan(&c.ID, &c.Title, &c.Year, &c.Artist, &c.PublishedBy, &c.Path); err != nil {
			return nil, err
		}
		covers = append(covers, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return covers, nil
}

func QueryAlbumCoverSummaryByAssetID(ctx context.Context, assetIDs []uuid.UUID) (map[uuid.UUID]AlbumCoverSummary, error) {
	if len(assetIDs) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(assetIDs))
	args := make([]any, len(assetIDs))
	for i, id := range assetIDs {
		placeholders[i] = "$" + strconv.Itoa(i+1)
		args[i] = id
	}

	rows, err := db.QueryContext(
		ctx,
		`SELECT DISTINCT
			ac.id, ac.title, ac.year, ac.artist, ac.published_by,asset.id, asset.path
		FROM album_cover ac
		INNER JOIN album_cover_asset aca ON aca.album_cover_id = ac.id
		INNER JOIN asset ON aca.asset_id = asset.id
		WHERE aca.asset_id IN (`+strings.Join(placeholders, ", ")+`)`,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	covers := make(map[uuid.UUID]AlbumCoverSummary, len(assetIDs))
	for rows.Next() {
		var ac AlbumCoverSummary
		var assetID uuid.UUID
		if err := rows.Scan(&ac.ID, &ac.Title, &ac.Year, &ac.Artist, &ac.PublishedBy, &assetID, &ac.Path); err != nil {
			return nil, err
		}
		covers[assetID] = ac
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return covers, nil
}

func InsertAlbumCoverTx(ctx context.Context, tx *sql.Tx, ac AlbumCover) error {
	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO album_cover(id, title, artist, rhythms, representation_technique, predominance, colors, title_syntax, published_by)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		ac.ID, ac.Title, ac.Artist, ac.Rhythms, ac.ReprTechnique, ac.Predominance, ac.Colors, ac.TitleSyntax, ac.PublishedBy,
	)

	return err
}

func InsertAlbumCoverAssetTx(ctx context.Context, tx *sql.Tx, albumCoverID uuid.UUID, assetID uuid.UUID) error {
	_, err := tx.ExecContext(ctx, `INSERT INTO album_cover_asset(album_cover_id, asset_id) VALUES ($1, $2)`, albumCoverID, assetID)
	return err
}
