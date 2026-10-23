import { HttpClient, HttpParams } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { map, Observable, tap } from 'rxjs';
import { AssetSearchParams, AssetSearchResponse } from '../interfaces/asset-search.interface';
import { AssetResponse } from '../interfaces/asset-response.interface';
import { EnvironmentConfigService } from '../../shared/services/environment-config.service';
import { Asset } from '../interfaces/asset.interface';
import { AssetMapper } from '../mappers/asset.mapper';

export type AssetKind = 'image' | 'audio';
export type EmbeddingModel = 'jina-clip-v2';

@Injectable({
  providedIn: 'root',
})
export class AssetsService {
  private http = inject(HttpClient);
  private env = inject(EnvironmentConfigService);
  private apiUrl = this.env.get(EnvironmentConfigService.API_URL);

  upload(asset: File): Observable<Asset> {
    const body = new FormData();
    body.append('file', asset, asset.name);

    return this.http
      .post<AssetResponse>(`${this.apiUrl}/assets`, body)
      .pipe(map(AssetMapper.fromResponse));
  }

  search({ q, kind, model, limit }: AssetSearchParams): Observable<AssetSearchResponse> {
    let params = new HttpParams().set('q', q).set('kind', kind);
    if (model) params = params.set('model', model);
    if (limit !== undefined) params = params.set('limit', limit);

    return this.http.get<AssetSearchResponse>(`${this.apiUrl}/assets/search`, { params });
  }
}
