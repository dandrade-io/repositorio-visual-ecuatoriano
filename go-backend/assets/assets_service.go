package assets

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"path/filepath"
	"repository/database"
	"repository/embedding"
	"repository/files"
	"repository/media"
	"uuid"

	"github.com/pgvector/pgvector-go"
)

type Service struct{}

type CreateInput struct {
	File   io.Reader
	Kind   string
	UserID uuid.UUID
	Embed  bool
}

func (s *Service) Create(ctx context.Context, in CreateInput) (database.Asset, error) {
	id := uuid.NewV7()
	path, err := s.SaveLocally(in.File, id.String(), in.Kind)
	if err != nil {
		return database.Asset{}, fmt.Errorf("unable to save asset locally: %w", err)
	}

	var asset = database.Asset{ID: id, Path: path, Kind: in.Kind, UploadedBy: in.UserID}

	if !in.Embed {
		if err := database.InsertAsset(ctx, asset); err != nil {
			return database.Asset{}, fmt.Errorf("failed to insert asset: %w", err)
		}
		return asset, nil
	}

	embeddings, err := embedding.EmbedAsset(ctx, asset.Path, asset.Kind, []embedding.Model{embedding.JinaClipV2, embedding.SigLIP2})
	if err != nil {
		return database.Asset{}, fmt.Errorf("failed to generate embeddings: %w", err)
	}

	err = database.WithTx(ctx, func(tx *sql.Tx) error {
		if err := database.InsertAssetTx(ctx, tx, asset); err != nil {
			return err
		}

		for _, e := range embeddings {
			entity := database.Embedding{AssetID: asset.ID, Vector: pgvector.NewVector(e.Embedding)}
			if err := database.InsertEmbeddingTx(ctx, tx, entity, e.Model); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return database.Asset{}, fmt.Errorf("failed to insert asset: %w", err)
	}

	return asset, nil
}

// SaveLocally saves file in the filesystem and returns the relative path
func (s *Service) SaveLocally(file io.Reader, filename string, kind string) (string, error) {
	switch kind {
	case "image":
		compressed, err := media.CompressImg(file, 85)
		if err != nil {
			return "", err
		}
		file = bytes.NewReader(compressed)
	}

	path, err := files.SaveFile(file, filepath.Join(kind, filename))
	if err != nil {
		return "", err
	}

	return path, nil
}

func (s *Service) EmbedByAssetID(ctx context.Context, assetID uuid.UUID, model embedding.Model) error {
	asset, err := database.GetAssetByID(ctx, assetID)
	if err != nil {
		return err
	}

	embeddings, err := embedding.EmbedAsset(ctx, asset.Path, asset.Kind, []embedding.Model{model})
	if err != nil || len(embeddings) == 0 {
		return fmt.Errorf("failed to generate embedding of asset %s: %w", assetID, err)
	}

	entity := database.Embedding{AssetID: asset.ID, Vector: pgvector.NewVector(embeddings[0].Embedding)}
	if err := database.InsertEmbedding(ctx, entity, model); err != nil {
		return fmt.Errorf("failed to store embeddings of asset %s in db: %w", assetID, err)
	}

	return nil
}

func (s *Service) Search(ctx context.Context, query string, m embedding.Model, limit int) (searchResults, error) {
	embedding, err := embedding.EmbedText(ctx, query, m)
	if err != nil {
		return searchResults{}, fmt.Errorf("unable to compute embedding for text: %w", err)
	}

	kBest, err := database.GetKBest(ctx, m, pgvector.NewVector(embedding), limit)
	if err != nil {
		return searchResults{}, fmt.Errorf("unable to search for k-best coincidence: %w", err)
	}

	assetIDs := make([]uuid.UUID, len(kBest))
	for i := range kBest {
		assetIDs[i] = kBest[i].AssetId
	}

	albumCovers, err := database.QueryAlbumCoverSummaryByAssetID(ctx, assetIDs)
	if err != nil {
		return searchResults{}, fmt.Errorf("unable to query album covers by asset ids: %w", err)
	}

	sr := searchResults{Results: make([]result, 0, len(albumCovers))}
	for _, b := range kBest {
		r := result{Score: b.Similarity, Type: Assets}

		if ac, ok := albumCovers[b.AssetId]; ok {
			r.Content = ac
			r.Type = AlbumCover
		} else {
			r.Content = map[string]any{"path": "image/" + b.AssetId.String()}
		}

		sr.Results = append(sr.Results, r)
	}

	return sr, nil
}
