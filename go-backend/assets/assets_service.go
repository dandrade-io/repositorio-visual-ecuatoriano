package assets

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"path/filepath"
	"repository/config"
	"repository/database"
	embedservice "repository/embed"
	"repository/files"
	"uuid"

	"github.com/pgvector/pgvector-go"
)

type Service struct{}

func (s *Service) Ingest(ctx context.Context, file io.Reader, kind string, model database.Model) error {
	asset, err := s.Create(file, kind)
	if err != nil {
		return err
	}

	if err := s.embed(ctx, asset, model); err != nil {
		return err
	}

	return nil
}

func (s *Service) Create(file io.Reader, kind string) (database.Asset, error) {
	if kind == "image" {
		compressed, err := compressImg(file, 85)
		if err != nil {
			return database.Asset{}, err
		}
		file = bytes.NewReader(compressed)
	}

	id := uuid.NewV7()
	path, err := files.SaveFile(file, filepath.Join(kind, id.String()))
	if err != nil {
		return database.Asset{}, fmt.Errorf("unable to save asset locally: %w", err)
	}

	if err := database.CreateAsset(database.Asset{ID: id, Path: path, Kind: "image"}); err != nil {
		return database.Asset{}, fmt.Errorf("unable to store asset in db: %w", err)
	}

	return database.Asset{ID: id, Path: path, Kind: kind}, nil
}

func (s *Service) EmbedByAssetID(ctx context.Context, assetID uuid.UUID, model database.Model) error {
	asset, err := database.GetAssetByID(ctx, assetID)
	if err != nil {
		return err
	}

	return s.embed(ctx, asset, model)
}

func (s *Service) embed(ctx context.Context, asset database.Asset, model database.Model) error {
	body := &bytes.Buffer{}
	payload := extractEmbeddingsReq{AssetURL: config.PublicBaseURL + "/api/v1/assets/" + asset.Path, Model: model.String()}
	if err := json.NewEncoder(body).Encode(&payload); err != nil {
		return err
	}

	url := config.EmbedServiceURL + "/api/v1/embed/" + asset.Kind
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make http petition: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to make request (%d): %s", resp.StatusCode, string(b))
	}

	var embeddingResp extractEmbeddingResp
	if err := json.NewDecoder(resp.Body).Decode(&embeddingResp); err != nil {
		return err
	}

	entity := database.Embedding{AssetID: asset.ID, Vector: pgvector.NewVector(embeddingResp.Embedding)}
	if err := database.CreateEmbedding(entity, model); err != nil {
		return fmt.Errorf("unable to store embeddings in db: %w", err)
	}

	return nil
}

func (s *Service) search(ctx context.Context, query string, model database.Model, limit int) ([]database.EmbeddingCoincidence, error) {
	embed, err := embedservice.Text(ctx, query, model)
	if err != nil {
		return nil, fmt.Errorf("unable to compute embedding for text: %w", err)
	}

	best, err := database.GetKBest(ctx, model, pgvector.NewVector(embed.Embedding), limit)
	if err != nil {
		return nil, fmt.Errorf("unable to search for k-best coincidence: %w", err)
	}

	return best, nil
}

func compressImg(r io.Reader, quality int) ([]byte, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality}); err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	return buf.Bytes(), nil
}
