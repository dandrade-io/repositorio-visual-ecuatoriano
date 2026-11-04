import { Predominance } from '../enums/prodominance.enum';
import { Technique } from '../enums/technique.enum';
import { TitleSyntax } from '../enums/title-syntax.enum';

export interface AlbumCoverMetadataResponse {
  title: string | null;
  year: number | null;
  artist: string | null;
  rhythms: string[];
  representationTechnique: string;
  predominance: string;
  dominantColors: string[];
  titleSyntax: string;
}

export interface AlbumCoverMetadata {
  title: string | null;
  year: number | null;
  artist: string | null;
  rhythms: string[];
  representationTechnique: Technique | null;
  predominance: Predominance | null;
  dominantColors: string[];
  titleSyntax: TitleSyntax | null;
}
