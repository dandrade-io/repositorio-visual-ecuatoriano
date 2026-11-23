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
	File     io.Reader
	Modality embedding.Modality
	UserID   uuid.UUID
	Embed    bool
}

func (s *Service) Create(ctx context.Context, in CreateInput) (database.Asset, error) {
	id := uuid.NewV7()
	path, err := s.SaveLocally(in.File, id.String(), in.Modality)
	if err != nil {
		return database.Asset{}, fmt.Errorf("unable to save asset locally: %w", err)
	}

	var asset = database.Asset{ID: id, Path: path, Kind: in.Modality.String(), UploadedBy: in.UserID}

	if !in.Embed {
		if err := database.InsertAsset(ctx, asset); err != nil {
			return database.Asset{}, fmt.Errorf("failed to insert asset: %w", err)
		}
		return asset, nil
	}

	var models []embedding.Model
	switch in.Modality {
	case embedding.Image:
		models = []embedding.Model{embedding.JinaClipV2, embedding.SigLIP2}
	case embedding.Audio:
		models = []embedding.Model{embedding.CLaMP3}
	}

	embeddings, err := embedding.EmbedAssetURL(ctx, asset.Path, embedding.Modality(asset.Kind), models)
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
func (s *Service) SaveLocally(file io.Reader, filename string, modality embedding.Modality) (string, error) {
	switch modality {
	case embedding.Image:
		compressed, err := media.CompressImg(file, 85)
		if err != nil {
			return "", err
		}
		file = bytes.NewReader(compressed)
	}

	path, err := files.SaveFile(file, filepath.Join(modality.String(), filename))
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

	embeddings, err := embedding.EmbedAssetURL(ctx, asset.Path, embedding.Modality(asset.Kind), []embedding.Model{model})
	if err != nil || len(embeddings) == 0 {
		return fmt.Errorf("failed to generate embedding of asset %s: %w", assetID, err)
	}

	entity := database.Embedding{AssetID: asset.ID, Vector: pgvector.NewVector(embeddings[0].Embedding)}
	if err := database.InsertEmbedding(ctx, entity, model); err != nil {
		return fmt.Errorf("failed to store embeddings of asset %s in db: %w", assetID, err)
	}

	return nil
}

func (s *Service) SearchByText(ctx context.Context, query string, model embedding.Model, modality embedding.Modality, limit int) (searchResult, error) {
	embedding, err := embedding.EmbedText(ctx, query, model)
	if err != nil {
		return searchResult{}, fmt.Errorf("unable to compute embedding for text: %w", err)
	}

	return s.search(ctx, embedding, model, limit)
}

func (s *Service) SearchByAsset(ctx context.Context, asset io.Reader, model embedding.Model, modality embedding.Modality, limit int) (searchResult, error) {
	embedding, err := embedding.EmbedAsset(ctx, asset, model, modality)
	if err != nil {
		return searchResult{}, fmt.Errorf("unable to compute embedding for asset: %w", err)
	}

	return s.search(ctx, embedding, model, limit)
}

func (s *Service) search(ctx context.Context, embedding []float32, model embedding.Model, limit int) (searchResult, error) {
	kBest, err := database.GetKBest(ctx, model, pgvector.NewVector(embedding), limit)
	if err != nil {
		return searchResult{}, fmt.Errorf("unable to search for k-best coincidence: %w", err)
	}

	assetIDs := make([]uuid.UUID, len(kBest))
	for i := range kBest {
		assetIDs[i] = kBest[i].AssetId
	}

	albumCovers, err := database.QueryAlbumCoverSummaryByAssetID(ctx, assetIDs)
	if err != nil {
		return searchResult{}, fmt.Errorf("unable to query album covers by asset ids: %w", err)
	}

	sr := searchResult{Results: make([]result, 0, len(albumCovers))}
	for _, b := range kBest {
		r := result{Score: b.Similarity, Type: Assets}

		if ac, ok := albumCovers[b.AssetId]; ok {
			r.Content = ac
			r.Type = AlbumCover
		} else {
			r.Content = map[string]any{"id": b.AssetId, "kind": b.Kind, "path": b.Path}
		}

		sr.Results = append(sr.Results, r)
	}

	return sr, nil
}
