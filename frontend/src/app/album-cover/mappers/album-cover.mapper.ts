import { parseEnum } from '../../shared/parse-enum';
import { Predominance } from '../enums/prodominance.enum';
import { Technique } from '../enums/technique.enum';
import { TitleSyntax } from '../enums/title-syntax.enum';
import { AlbumCover, AlbumCoverResponse } from '../interfaces/album-cover-metadata-response';

export class AlbumCoverMapper {
  static fromResponse(resp: AlbumCoverResponse): AlbumCover {
    console.log(resp);

    return {
      title: resp.title,
      year: resp.year,
      artist: resp.artist,
      rhythms: resp.rhythms,
      representationTechnique: parseEnum(Technique, resp.representationTechnique),
      graphicPredominance: parseEnum(Predominance, resp.graphicPredominance),
      dominantColors: resp.dominantColors,
      titleSyntax: parseEnum(TitleSyntax, resp.titleSyntax),
    };
  }
}
