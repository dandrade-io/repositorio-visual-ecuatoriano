export interface AsssetSearchResponse {
  results: SearchResult[];
}

export interface SearchResult {
  assetId: string;
  similarity: number;
}

// Raw shape returned by the Go backend (snake_case).
export interface AssetSearchApiResponse {
  results: { asset_id: string; similarity: number }[];
  num_results: number;
}
