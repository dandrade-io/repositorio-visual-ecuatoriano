package database

import (
	"context"
	"database/sql"
	"fmt"
	"repository/embedding"
	"uuid"

	"github.com/pgvector/pgvector-go"
)

func embeddingModelToTableName(m embedding.Model) string {
	switch m {
	case embedding.JinaClipV2:
		return "jina_clip_v2"
	case embedding.SigLIP2:
		return "siglip2"
	case embedding.CLaMP3:
		return "clamp3"
	default:
		return ""
	}
}

type Embedding struct {
	AssetID uuid.UUID
	Vector  pgvector.Vector
}

func createEmbeddingQuery(e Embedding, m embedding.Model) (string, []any) {
	return "INSERT INTO embedding_" + embeddingModelToTableName(m) + " (asset_id, embedding) VALUES($1, $2)", []any{e.AssetID, e.Vector}
}

func InsertEmbedding(ctx context.Context, e Embedding, m embedding.Model) error {
	q, args := createEmbeddingQuery(e, m)
	_, err := db.ExecContext(ctx, q, args...)

	if err != nil {
		return err
	}

	return nil
}

func InsertEmbeddingTx(ctx context.Context, tx *sql.Tx, e Embedding, m embedding.Model) error {
	q, args := createEmbeddingQuery(e, m)
	_, err := tx.ExecContext(ctx, q, args...)

	return err
}

type EmbeddingCoincidence struct {
	AssetId    uuid.UUID
	Kind       string
	Path       string
	Similarity float32
}

func GetKBest(ctx context.Context, m embedding.Model, vector pgvector.Vector, limit int) ([]EmbeddingCoincidence, error) {
	query := `
		WITH knn AS (
			SELECT asset_id, embedding <=> $1 AS dist
			FROM embedding_` + embeddingModelToTableName(m) + `
			ORDER BY embedding <=> $1
			LIMIT $2
		)
		SELECT knn.asset_id, a.kind, a.path, 1 - knn.dist AS similarity
		FROM knn
		JOIN asset a ON a.id = knn.asset_id
		ORDER BY knn.dist;
		`

	rows, err := db.QueryContext(ctx, query, vector, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]EmbeddingCoincidence, 0, limit)

	for rows.Next() {
		var e EmbeddingCoincidence
		if err := rows.Scan(&e.AssetId, &e.Kind, &e.Path, &e.Similarity); err != nil {
			return nil, err
		}
		results = append(results, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return results, nil

}
