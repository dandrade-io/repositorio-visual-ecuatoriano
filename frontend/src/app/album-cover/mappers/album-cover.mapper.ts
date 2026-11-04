import { isVerifierRole } from '../../auth/services/auth.service';
import { parseEnum } from '../../shared/parse-enum';
import { Predominance } from '../enums/prodominance.enum';
import { Technique } from '../enums/technique.enum';
import { TitleSyntax } from '../enums/title-syntax.enum';
import { isVerifiableField } from '../enums/verifiable_fields';
import { AlbumCoverMetadata, AlbumCoverMetadataResponse } from '../interfaces/album-cover-metadata';
import {
  AlbumCover,
  AlbumCoverCheck,
  AlbumCoverCheckResponse,
  AlbumCoverResponse,
} from '../interfaces/album-cover.interface';

export class AlbumCoverMapper {
  static fromMetadataResponse(response: AlbumCoverMetadataResponse): AlbumCoverMetadata {
    return {
      title: response.title,
      year: response.year,
      artist: response.artist,
      rhythms: response.rhythms,
      representationTechnique: parseEnum(Technique, response.representationTechnique),
      predominance: parseEnum(Predominance, response.predominance),
      dominantColors: response.dominantColors,
      titleSyntax: parseEnum(TitleSyntax, response.titleSyntax),
    };
  }

  static fromResponse(response: AlbumCoverResponse): AlbumCover {
    return {
      id: response.id,
      path: response.path,
      checks: response.checks.map(AlbumCoverMapper.toCheck).filter((c) => c !== null),
      ...AlbumCoverMapper.fromMetadataResponse(response),
    };
  }

  private static toCheck(check: AlbumCoverCheckResponse): AlbumCoverCheck | null {
    const fieldName = snakeToCamelCase(check.fieldName);
    if (!isVerifiableField(fieldName) || !isVerifierRole(check.verifierRole)) return null;
    return { ...check, fieldName, verifierRole: check.verifierRole };
  }
}

function snakeToCamelCase(str: string): string {
  return str.replace(/_([a-z])/g, (_, letter: string) => letter.toUpperCase());
}
