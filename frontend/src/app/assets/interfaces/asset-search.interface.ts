import { AssetKind, EmbeddingModel } from '../services/assets.service';

export interface AssetSearchParams {
  q: string;
  kind: AssetKind;
  model?: EmbeddingModel;
  limit?: number;
}

export interface AssetSearchResponse {
  results: Result[];
}

export interface Result {
  score: number;
  type: string;
  content: Content;
}

export interface Content {
  id: string;
  title: string | null;
  year: number | null;
  artist: string | null;
  published_by: string;
  path: string;
}
