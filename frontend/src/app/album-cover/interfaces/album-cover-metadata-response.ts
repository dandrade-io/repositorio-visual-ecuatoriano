import { Predominance } from '../enums/prodominance.enum';
import { Technique } from '../enums/technique.enum';
import { TitleSyntax } from '../enums/title-syntax.enum';

export interface AlbumCoverResponse {
  title: string | null;
  year: number | null;
  artist: string | null;
  rhythms: string[];
  representationTechnique: string;
  graphicPredominance: string;
  dominantColors: string[];
  titleSyntax: string;
}

export interface AlbumCover {
  title: string | null;
  year: number | null;
  artist: string | null;
  rhythms: string[];
  representationTechnique: Technique | null;
  graphicPredominance: Predominance | null;
  dominantColors: string[];
  titleSyntax: TitleSyntax | null;
}
