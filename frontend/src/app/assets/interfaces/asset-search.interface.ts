import { ArtworkType } from '../../artwork/interfaces/artwork-base.interface';
import { Modality, EmbeddingModel } from '../services/assets.service';

export interface AssetSearchParams {
  q: string;
  kind: Modality;
  model?: EmbeddingModel;
  limit?: number;
}

export interface AssetSearchAssetParams {
  asset: File;
  modality: Modality;
  model: EmbeddingModel;
  limit: number;
}

export interface AssetSearchResponse {
  results: Result[];
}

export interface Result {
  score: number;
  type: ArtworkType;
  content: Content;
}

export interface Content {
  id: string;
  title: string | null;
  year: number | null;
  artist: string | null;
  published_by: string;
  kind: string | null;
  path: string;
}
