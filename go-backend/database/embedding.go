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

func InsertEmbedding(e Embedding, m embedding.Model) error {
	q, args := createEmbeddingQuery(e, m)
	_, err := db.Exec(q, args...)

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
	Similarity float32
}

func GetKBest(ctx context.Context, m embedding.Model, vector pgvector.Vector, limit int) ([]EmbeddingCoincidence, error) {
	query := `
		SELECT asset_id, (1 - (embedding <=> $1)) AS similarity
		FROM embedding_` + embeddingModelToTableName(m) + `
		ORDER BY embedding <=> $1
		LIMIT $2
		`

	rows, err := db.QueryContext(ctx, query, vector, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]EmbeddingCoincidence, 0, limit)

	for rows.Next() {
		var ec EmbeddingCoincidence
		if err := rows.Scan(&ec.AssetId, &ec.Similarity); err != nil {
			return nil, err
		}
		results = append(results, ec)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return results, nil

}
