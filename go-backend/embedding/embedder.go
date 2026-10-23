package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"repository/config"
)

func EmbedText(ctx context.Context, text string, model Model) ([]float32, error) {
	type embedTextReq struct {
		Text  string `json:"text"`
		Model string `json:"model"`
	}

	type embedTextResp struct {
		Embedding []float32 `json:"embedding"`
	}

	body := &bytes.Buffer{}
	payload := &embedTextReq{Text: text, Model: model.String()}
	if err := json.NewEncoder(body).Encode(payload); err != nil {
		return nil, err
	}

	url := config.EmbedServiceURL + "/api/v1/embed/text"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make http petition: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to make request (%d): %s", resp.StatusCode, string(b))
	}

	var embedResp embedTextResp
	if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
		return nil, err
	}

	return embedResp.Embedding, nil
}

func EmbedAsset(ctx context.Context, assetPath string, assetKind string, model Model) ([]float32, error) {
	type extractEmbeddingReq struct {
		AssetURL string `json:"asset_url"`
		Model    string `json:"model"`
	}

	type extractEmbeddingResp struct {
		Model     string    `json:"model"`
		Embedding []float32 `json:"embedding"`
	}

	body := &bytes.Buffer{}
	payload := extractEmbeddingReq{AssetURL: config.PublicBaseURL + "/api/v1/assets/" + assetPath, Model: model.String()}
	if err := json.NewEncoder(body).Encode(&payload); err != nil {
		return nil, err
	}

	url := config.EmbedServiceURL + "/api/v1/embed/" + assetKind
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make http petition: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to make request (%d): %s", resp.StatusCode, string(b))
	}

	var embeddingResp extractEmbeddingResp
	if err := json.NewDecoder(resp.Body).Decode(&embeddingResp); err != nil {
		return nil, err
	}

	return embeddingResp.Embedding, nil
}
