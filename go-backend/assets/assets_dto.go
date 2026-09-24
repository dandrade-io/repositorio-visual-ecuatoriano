package assets

import "uuid"

type assetResp struct {
	ID   uuid.UUID `json:"id"`
	Path string    `json:"path"`
}

type embedReq struct {
	Model string `json:"model"`
}

type extractEmbeddingsReq struct {
	AssetURL string `json:"asset_url"`
	Model    string `json:"model"`
}

type extractEmbeddingResp struct {
	Model     string    `json:"model"`
	Embedding []float32 `json:"embedding"`
}

type searchResp struct {
	Results    []searchResult `json:"results"`
	NumResults int            `json:"num_results"`
}

type searchResult struct {
	AssetId    uuid.UUID `json:"asset_id"`
	Similarity float32   `json:"similarity"`
}
