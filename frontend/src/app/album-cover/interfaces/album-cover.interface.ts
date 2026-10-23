import { AlbumCoverMetadata, AlbumCoverMetadataResponse } from './album-cover-metadata';

export interface AlbumCoverResponse extends AlbumCoverMetadataResponse {
  id: string;
}

export interface AlbumCover extends AlbumCoverMetadata {
  id: string;
}
