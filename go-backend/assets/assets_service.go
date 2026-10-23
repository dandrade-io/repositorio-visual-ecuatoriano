package assets

import (
	"bytes"
	"context"
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

func (s *Service) Ingest(ctx context.Context, file io.Reader, kind string, model embedding.Model, userID uuid.UUID) error {
	asset, err := s.Create(ctx, file, kind, userID)
	if err != nil {
		return err
	}

	if err := s.embedAndSave(ctx, asset, model); err != nil {
		return err
	}

	return nil
}

func (s *Service) Create(ctx context.Context, file io.Reader, kind string, userID uuid.UUID) (database.Asset, error) {
	id, path, err := s.SaveLocally(file, kind)
	if err != nil {
		return database.Asset{}, fmt.Errorf("unable to save asset locally: %w", err)
	}

	var a = database.Asset{ID: id, Path: path, Kind: kind, UploadedBy: userID}
	if err := database.InsertAsset(ctx, a); err != nil {
		return database.Asset{}, fmt.Errorf("failed to insert asset: %w", err)
	}

	return a, nil
}

// SaveLocally saves file in the filesystem
func (s *Service) SaveLocally(file io.Reader, kind string) (uuid.UUID, string, error) {
	switch kind {
	case "image":
		compressed, err := media.CompressImg(file, 85)
		if err != nil {
			return uuid.Nil(), "", err
		}
		file = bytes.NewReader(compressed)
	}

	id := uuid.NewV7()
	path, err := files.SaveFile(file, filepath.Join(kind, id.String()))
	if err != nil {
		return uuid.Nil(), "", err
	}

	return id, path, nil
}

func (s *Service) EmbedByAssetID(ctx context.Context, assetID uuid.UUID, model embedding.Model) error {
	asset, err := database.GetAssetByID(ctx, assetID)
	if err != nil {
		return err
	}

	return s.embedAndSave(ctx, asset, model)
}

func (s *Service) embedAndSave(ctx context.Context, asset database.Asset, model embedding.Model) error {
	embedding, err := embedding.EmbedAsset(ctx, asset.Path, asset.Kind, model)
	if err != nil {
		return fmt.Errorf("failed to embed: %w", err)
	}

	entity := database.Embedding{AssetID: asset.ID, Vector: pgvector.NewVector(embedding)}
	if err := database.InsertEmbedding(entity, model); err != nil {
		return fmt.Errorf("unable to store embeddings in db: %w", err)
	}

	return nil
}

func (s *Service) search(ctx context.Context, query string, m embedding.Model, limit int) (searchResults, error) {
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

	albumCovers, err := database.QueryAlbumCoverByAssetID(ctx, assetIDs)
	if err != nil {
		return searchResults{}, fmt.Errorf("unable to query album covers by asset ids: %w", err)
	}

	sr := searchResults{Results: make([]result, 0, len(albumCovers))}
	for _, b := range kBest {
		if ac, ok := albumCovers[b.AssetId]; ok {
			r := result{Score: b.Similarity, Type: AlbumCover, Content: ac}
			sr.Results = append(sr.Results, r)
		}
	}

	return sr, nil
}
