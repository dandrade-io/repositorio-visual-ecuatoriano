package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"uuid"
)

type AlbumCoverField string

const (
	Title                   = "title"
	Year                    = "year"
	Artist                  = "artist"
	Rhythms                 = "rhythms"
	RepresentationTechnique = "representation_technique"
	Predominance            = "predominance"
	Colors                  = "colors"
	TitleSyntax             = "title_syntax"
)

func (f AlbumCoverField) Valid() bool {
	switch f {
	case Title, Year, Artist, Rhythms, RepresentationTechnique, Predominance, Colors, TitleSyntax:
		return true
	}

	return false
}

func (f AlbumCoverField) Value() string {
	return string(f)
}

func InsertAlbumCoverCheckTx(
	ctx context.Context, tx *sql.Tx,
	fields map[AlbumCoverField]struct{}, albumCoverID uuid.UUID, verifierID uuid.UUID, verifierRol string,
) error {
	args := make([]any, 0, len(fields)*4)
	placeholders := make([]string, 0, len(fields))

	i := 0
	for f := range fields {
		if !f.Valid() {
			return fmt.Errorf("invalid field: %s", f.Value())
		}

		b := (i * 4) + 1
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d,$%d)", b, b+1, b+2, b+3))
		args = append(args, f.Value(), albumCoverID, verifierID, verifierRol)
		i++
	}

	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO album_cover_check (field_name, album_cover_id, verifier_id, verifier_role)
		VALUES `+strings.Join(placeholders, ","),
		args...,
	)

	return err
}

type AlbumCoverCheck struct {
	FieldName    string    `json:"fieldName"`
	VerifierID   uuid.UUID `json:"verifierId"`
	VerifierRole string    `json:"verifierRole"`
}

func QueryAlbumCoverChecksByAlbumCoverID(ctx context.Context, albumCoverID uuid.UUID) ([]AlbumCoverCheck, error) {
	rows, err := db.QueryContext(ctx,
		`SELECT field_name, verifier_id, verifier_role FROM album_cover_check WHERE album_cover_id = $1`,
		albumCoverID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	checks := make([]AlbumCoverCheck, 0, 8)
	for rows.Next() {
		var c AlbumCoverCheck
		if err := rows.Scan(&c.FieldName, &c.VerifierID, &c.VerifierRole); err != nil {
			return nil, err
		}
		checks = append(checks, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return checks, nil
}
