import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { AlbumCoverMetadata, AlbumCoverMetadataResponse } from '../interfaces/album-cover-metadata';
import { map, Observable } from 'rxjs';
import { AlbumCoverMapper } from '../mappers/album-cover.mapper';
import { EnvironmentConfigService } from '../../shared/services/environment-config.service';
import { AlbumCover, AlbumCoverResponse } from '../interfaces/album-cover.interface';
import { AlbumCoverSummaryResponse } from '../interfaces/album-cover-summary.interface';

@Injectable({ providedIn: 'root' })
export class AlbumCoverService {
  private http = inject(HttpClient);
  private env = inject(EnvironmentConfigService);
  private apiUrl = this.env.get(EnvironmentConfigService.API_URL);

  fetchSummary(): Observable<AlbumCoverSummaryResponse> {
    return this.http.get<AlbumCoverSummaryResponse>(`${this.apiUrl}/album-cover`);
  }

  extractMetadata(assetId: string): Observable<AlbumCoverMetadata> {
    const url = `${this.apiUrl}/assets/${assetId}/metadata/album-cover`;

    return this.http
      .post<AlbumCoverMetadataResponse>(url, null)
      .pipe(map(AlbumCoverMapper.fromMetadataResponse));
  }

  upload(metadata: AlbumCoverMetadata, assetId: string): Observable<AlbumCover> {
    return this.http
      .post<AlbumCoverResponse>(`${this.apiUrl}/album-cover`, { ...metadata, assetId })
      .pipe(map(AlbumCoverMapper.fromResponse));
  }
}
