import { HttpClient, HttpParams } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { map, Observable, tap } from 'rxjs';
import {
  AssetSearchApiResponse,
  AsssetSearchResponse,
} from '../interfaces/asset-search-response.interface';

export type AssetKind = 'image' | 'audio';
export type EmbeddingModel = 'jina-clip-v2';

export interface AssetSearchParams {
  q: string;
  kind: AssetKind;
  model?: EmbeddingModel;
  limit?: number;
}

declare global {
  interface Window {
    __env?: { BACKEND_BASE_URL?: string };
  }
}

const BACKEND_BASE_URL = (window.__env?.BACKEND_BASE_URL ?? 'http://localhost:8080').replace(/\/$/, '');

@Injectable({
  providedIn: 'root',
})
export class AssetsService {
  private http = inject(HttpClient);

  imageUrl(assetId: string | number): string {
    return `${BACKEND_BASE_URL}/api/v1/assets/image/${assetId}`;
  }

  search({ q, kind, model, limit }: AssetSearchParams): Observable<AsssetSearchResponse> {
    let params = new HttpParams().set('q', q).set('kind', kind);
    if (model) params = params.set('model', model);
    if (limit !== undefined) params = params.set('limit', limit);

    console.log({ q, kind, model, limit });

    return this.http
      .get<AssetSearchApiResponse>(`${BACKEND_BASE_URL}/api/v1/assets/search`, { params })
      .pipe(
        tap((resp) => console.log(resp)),
        map((resp) => ({
          results: resp.results.map((r) => ({ assetId: r.asset_id, similarity: r.similarity })),
        })),
      );
  }
}
