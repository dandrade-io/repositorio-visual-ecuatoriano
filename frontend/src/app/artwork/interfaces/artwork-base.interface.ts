export interface ArtworkBase {
  id: string;
  title: string | null;
  artist?: string | null;
  year?: number | null;
  imageUrl: string;
  score: number | null;
}
