import { HttpClient } from '@angular/common/http';
import { inject, Service } from '@angular/core';
import { AlbumCoverMetadata, AlbumCoverMetadataResponse } from '../interfaces/album-cover-metadata';
import { map, Observable } from 'rxjs';
import { AlbumCoverMapper } from '../mappers/album-cover.mapper';
import { EnvironmentConfigService } from '../../shared/services/environment-config.service';
import { AlbumCover, AlbumCoverResponse } from '../interfaces/album-cover.interface';
import { AlbumCoverSummaryResponse } from '../interfaces/album-cover-summary.interface';
import { VerifiableField } from '../enums/verifiable_fields';

@Service()
export class AlbumCoverService {
  private http = inject(HttpClient);
  private env = inject(EnvironmentConfigService);
  private apiUrl = this.env.get(EnvironmentConfigService.API_URL);

  fetchSummary(): Observable<AlbumCoverSummaryResponse> {
    return this.http.get<AlbumCoverSummaryResponse>(`${this.apiUrl}/album-cover`);
  }

  fetchById(id: string): Observable<AlbumCover> {
    return this.http
      .get<AlbumCoverResponse>(`${this.apiUrl}/album-cover/${id}`)
      .pipe(map(AlbumCoverMapper.fromResponse));
  }

  extractMetadata(assetId: string): Observable<AlbumCoverMetadata> {
    const url = `${this.apiUrl}/assets/${assetId}/metadata/album-cover`;

    return this.http
      .post<AlbumCoverMetadataResponse>(url, null)
      .pipe(map(AlbumCoverMapper.fromMetadataResponse));
  }

  upload(
    metadata: AlbumCoverMetadata,
    verifiedFields: Set<VerifiableField>,
    assetId: string,
  ): Observable<AlbumCover> {
    const payload = {
      ...metadata,
      verifiedFields: Array.from(verifiedFields).map(camelToSnakeCase),
      assetId,
    };

    return this.http
      .post<AlbumCoverResponse>(`${this.apiUrl}/album-cover`, payload)
      .pipe(map(AlbumCoverMapper.fromResponse));
  }
}

function camelToSnakeCase(str: string): string {
  return str.replace(/[A-Z]/g, (letter) => `_${letter.toLowerCase()}`);
}
