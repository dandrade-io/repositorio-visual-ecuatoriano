package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"repository/config"
)

type EmbedResult struct {
	Model     Model
	Embedding []float32
}

func EmbedText(ctx context.Context, text string, model Model) ([]float32, error) {
	type embedReq struct {
		Text  string `json:"text"`
		Model string `json:"model"`
	}

	type embedResp struct {
		Embedding []float32 `json:"embedding"`
	}

	body := &bytes.Buffer{}
	payload := &embedReq{Text: text, Model: model.String()}
	if err := json.NewEncoder(body).Encode(payload); err != nil {
		return nil, err
	}

	url := config.EmbedServiceURL + "/api/v1/embed/text"
	slog.Info(url + "AAAAA")
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("X-API-Key", config.EmbedServiceApiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make http petition: %w", err)
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode >= 400 {
		b, _ := io.ReadAll(httpResp.Body)
		return nil, fmt.Errorf("failed to make request (%d): %s", httpResp.StatusCode, string(b))
	}

	var resp embedResp
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return resp.Embedding, nil
}

func EmbedAsset(ctx context.Context, assetPath string, assetKind string, models []Model) ([]EmbedResult, error) {
	type embedReq struct {
		AssetURL string   `json:"asset_url"`
		Models   []string `json:"models"`
	}

	type embedResp struct {
		Results []struct {
			Model     string    `json:"model"`
			Embedding []float32 `json:"embedding"`
		} `json:"results"`
	}

	if len(models) == 0 {
		return nil, fmt.Errorf("models must have a lenght greater than zero")
	}

	assetURL := config.PublicBaseURL + "/api/v1/assets/" + assetPath
	payload := embedReq{AssetURL: assetURL, Models: make([]string, len(models))}
	for i := range models {
		payload.Models[i] = models[i].String()
	}

	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(&payload); err != nil {
		return nil, err
	}

	url := config.EmbedServiceURL + "/api/v1/embed/" + assetKind
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("X-API-Key", config.EmbedServiceApiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make http petition: %w", err)
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode >= 400 {
		b, _ := io.ReadAll(httpResp.Body)
		return nil, fmt.Errorf("failed to make request (%d): %s", httpResp.StatusCode, string(b))
	}

	var resp embedResp
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, err
	}

	results := make([]EmbedResult, 0, len(resp.Results))
	for _, r := range resp.Results {
		m := Model(r.Model)
		if m.Valid() {
			continue
		}
		results = append(results, EmbedResult{Model: m, Embedding: r.Embedding})
	}

	return results, nil
}
