import { VerifierRole } from '../../auth/services/auth.service';
import { VerifiableField } from '../enums/verifiable_fields';
import { AlbumCoverMetadata, AlbumCoverMetadataResponse } from './album-cover-metadata';

export interface AlbumCoverResponse extends AlbumCoverMetadataResponse {
  id: string;
  path: string;
  checks: AlbumCoverCheckResponse[];
}

export interface AlbumCover extends AlbumCoverMetadata {
  id: string;
  path: string;
  checks: AlbumCoverCheck[];
}

export interface AlbumCoverCheckResponse {
  fieldName: string;
  verifierId: string;
  verifierRole: string;
}

export interface AlbumCoverCheck {
  fieldName: VerifiableField;
  verifierId: string;
  verifierRole: VerifierRole;
}
