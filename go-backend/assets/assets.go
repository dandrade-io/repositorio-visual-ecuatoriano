package assets

type AssetType string

const (
	AlbumCover AssetType = "album-cover"
)

type searchResults struct {
	Results []result `json:"results"`
}

type result struct {
	Score   float32   `json:"score"`
	Type    AssetType `json:"type"`
	Content any       `json:"content"`
}
