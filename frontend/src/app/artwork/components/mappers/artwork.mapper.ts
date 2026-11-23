import {
  AlbumCoverSummary,
  AlbumCoverSummaryResponse,
} from '../../../album-cover/interfaces/album-cover-summary.interface';
import { AssetSearchResponse } from '../../../assets/interfaces/asset-search.interface';
import { ArtworkBase } from '../../interfaces/artwork-base.interface';

export class ArtworkMapper {
  static fromAssetSearchResponse(response: AssetSearchResponse): ArtworkBase[] {
    return response.results.map((r) => ({
      id: r.content.id,
      title: r.content.title,
      artist: r.content.artist,
      year: r.content.year,
      url: r.content.path,
      score: r.score,
      type: r.type,
      kind: r.content.kind ?? '',
    }));
  }

  static fromAlbumCoverSummaryResponse(response: AlbumCoverSummaryResponse): ArtworkBase[] {
    return response.results.map(ArtworkMapper.fromAlbumCoverSummary);
  }

  static fromAlbumCoverSummary(albumCover: AlbumCoverSummary): ArtworkBase {
    return {
      id: albumCover.id,
      title: albumCover.title,
      artist: albumCover.artist,
      year: albumCover.year,
      url: albumCover.path,
      score: null,
      type: 'album-cover',
      kind: 'image',
    };
  }
}
