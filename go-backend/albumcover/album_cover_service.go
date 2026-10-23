package albumcover

import (
	"context"
	"database/sql"
	"fmt"
	"repository/assets"
	"repository/database"
	"repository/embedding"
	"uuid"

	"github.com/pgvector/pgvector-go"
)

type Service struct {
	assetsSvc assets.Service
}

type CreateInput struct {
	Title         string
	Year          *int
	Artist        string
	Rhythms       []string
	ReprTechnique string
	Predominance  string
	Colors        []string
	TitleSyntax   string
	AssetID       uuid.UUID
	PublishedBy   uuid.UUID
}

func (in *CreateInput) toAlbumCover() database.AlbumCover {
	return database.AlbumCover{
		Title:         in.Title,
		Year:          in.Year,
		Artist:        in.Artist,
		Rhythms:       in.Rhythms,
		ReprTechnique: in.ReprTechnique,
		Predominance:  in.Predominance,
		Colors:        in.Colors,
		TitleSyntax:   in.TitleSyntax,
		PublishedBy:   in.PublishedBy,
	}
}

func (s *Service) FetchSummary(ctx context.Context, limit int) ([]database.AlbumCoverSummary, error) {
	return database.QueryAlbumCoverSummary(ctx, limit)
}

func (s *Service) Create(ctx context.Context, in CreateInput) (database.AlbumCover, error) {
	asset, err := database.GetAssetByID(ctx, in.AssetID)
	if err != nil {
		return database.AlbumCover{}, fmt.Errorf("failed to retrieve asset %s: %w", asset.ID, err)
	}

	m := embedding.JinaClipV2
	embedding, err := embedding.EmbedAsset(ctx, asset.Path, asset.Kind, m)
	if err != nil {
		return database.AlbumCover{}, fmt.Errorf("failed to generate embedding: %w", err)
	}

	albumCover := in.toAlbumCover()
	albumCover.ID = uuid.New()

	err = database.WithTx(ctx, func(tx *sql.Tx) error {
		e := database.Embedding{AssetID: asset.ID, Vector: pgvector.NewVector(embedding)}
		if err := database.InsertEmbeddingTx(ctx, tx, e, m); err != nil {
			return fmt.Errorf("failed to insert embedding: %w", err)
		}

		if err := database.InsertAlbumCoverTx(ctx, tx, albumCover); err != nil {
			return fmt.Errorf("failed to insert album-cover: %w", err)
		}

		if err := database.InsertAlbumCoverAssetTx(ctx, tx, albumCover.ID, asset.ID); err != nil {
			return fmt.Errorf("failed to insert album-cover-assset: %w", err)
		}

		err = database.ClaimAssetTx(ctx, tx, asset.ID, in.PublishedBy)
		return err
	})

	if err != nil {
		return database.AlbumCover{}, err
	}

	return albumCover, nil
}
