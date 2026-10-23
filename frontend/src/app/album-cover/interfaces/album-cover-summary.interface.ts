export interface AlbumCoverSummaryResponse {
  results: AlbumCoverSummary[];
}

export interface AlbumCoverSummary {
  id: string;
  title: string | null;
  year: number | null;
  artist: string | null;
  publishedBy: string;
  path: string;
}
