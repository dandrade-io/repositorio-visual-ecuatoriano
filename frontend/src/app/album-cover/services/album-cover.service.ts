import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { API_BASE_URL } from '../../shared/config/app.constants';
import { AlbumCover, AlbumCoverResponse } from '../interfaces/album-cover-metadata-response';
import { map, Observable } from 'rxjs';
import { AlbumCoverMapper } from '../mappers/album-cover.mapper';

@Injectable({
  providedIn: 'root',
})
export class AlbumCoverService {
  private http = inject(HttpClient);

  extractMetadata(img: File): Observable<AlbumCover> {
    const body = new FormData();
    body.append('asset', img, img.name);

    return this.http
      .post<AlbumCoverResponse>(`${API_BASE_URL}/metadata/album-cover`, body)
      .pipe(map(AlbumCoverMapper.fromResponse));
  }
}
