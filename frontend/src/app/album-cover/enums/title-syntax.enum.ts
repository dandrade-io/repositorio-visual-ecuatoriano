export enum TitleSyntax {
  Song = 'song',
  Image = 'image',
  Artist = 'artist',
  MusicKind = 'music_type',
  Company = 'company',
  Feeling = 'sentiment',
}

export const TITLE_SYNTAX_LABELS: Record<TitleSyntax, string> = {
  [TitleSyntax.Song]: 'La canción',
  [TitleSyntax.Image]: 'La imagen',
  [TitleSyntax.Artist]: 'El artist',
  [TitleSyntax.MusicKind]: 'El tipo de música',
  [TitleSyntax.Company]: 'La empresa',
  [TitleSyntax.Feeling]: 'El sentimiento',
};
