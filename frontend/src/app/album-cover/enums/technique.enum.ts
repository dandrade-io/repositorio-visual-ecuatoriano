export enum Technique {
  Photography = 'photography',
  Illustration = 'illustration',
  Collage = 'collage',
  Abstract = 'abstract',
  Simbolic = 'symbolic',
  TypographicComposition = 'typographic_composition',
}

export const TECHNIQUE_LABELS: Record<Technique, string> = {
  [Technique.Photography]: 'Fotografía',
  [Technique.Illustration]: 'Ilustración',
  [Technique.Collage]: 'Collage',
  [Technique.Abstract]: 'Abstracta',
  [Technique.Simbolic]: 'Simbólica',
  [Technique.TypographicComposition]: 'Composición tipográfica',
};
