import { parseEnum } from '../../shared/parse-enum';
import { Predominance } from '../enums/prodominance.enum';
import { Technique } from '../enums/technique.enum';
import { TitleSyntax } from '../enums/title-syntax.enum';
import { AlbumCoverMetadata, AlbumCoverMetadataResponse } from '../interfaces/album-cover-metadata';
import { AlbumCover, AlbumCoverResponse } from '../interfaces/album-cover.interface';

export class AlbumCoverMapper {
  static fromMetadataResponse(response: AlbumCoverMetadataResponse): AlbumCoverMetadata {
    return {
      title: response.title,
      year: response.year,
      artist: response.artist,
      rhythms: response.rhythms,
      representationTechnique: parseEnum(Technique, response.representationTechnique),
      graphicPredominance: parseEnum(Predominance, response.graphicPredominance),
      dominantColors: response.dominantColors,
      titleSyntax: parseEnum(TitleSyntax, response.titleSyntax),
    };
  }

  static fromResponse(response: AlbumCoverResponse): AlbumCover {
    return {
      id: response.id,
      ...AlbumCoverMapper.fromMetadataResponse(response),
    };
  }
}
