package database

import (
	"context"
	"fmt"
	"uuid"

	"github.com/pgvector/pgvector-go"
)

type Model int

const (
	JinaClipV2 Model = iota
)

func (m Model) String() string {
	switch m {
	case JinaClipV2:
		return "jina-clip-v2"
	default:
		return ""
	}
}

func (m Model) toTableName() string {
	switch m {
	case JinaClipV2:
		return "jina_clip_v2"
	default:
		return ""
	}
}

func ParseModel(val string) (Model, error) {
	switch val {
	case "jina-clip-v2":
		return JinaClipV2, nil
	default:
		return 0, fmt.Errorf("unknown model: %s", val)
	}
}

type Embedding struct {
	AssetID uuid.UUID
	Vector  pgvector.Vector
}

func CreateEmbedding(e Embedding, m Model) error {
	_, err := db.Exec("INSERT INTO embedding_"+m.toTableName()+" (asset_id, embedding) VALUES($1, $2)", e.AssetID, e.Vector)

	if err != nil {
		return err
	}

	return nil
}

type EmbeddingCoincidence struct {
	AssetId    uuid.UUID
	Similarity float32
}

func GetKBest(ctx context.Context, m Model, vector pgvector.Vector, limit int) ([]EmbeddingCoincidence, error) {
	query := `
		SELECT asset_id, (1 - (embedding <=> $1)) AS similarity
		FROM embedding_` + m.toTableName() + `
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

	return results, nil

}
