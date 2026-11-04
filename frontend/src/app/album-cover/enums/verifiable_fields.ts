export const VERIFIABLE_FIELDS = [
  'title',
  'year',
  'artist',
  'rhythms',
  'representationTechnique',
  'predominance',
  'titleSyntax',
] as const;

export type VerifiableField = (typeof VERIFIABLE_FIELDS)[number];

export function isVerifiableField(value: string): value is VerifiableField {
  return (VERIFIABLE_FIELDS as readonly string[]).includes(value);
}
