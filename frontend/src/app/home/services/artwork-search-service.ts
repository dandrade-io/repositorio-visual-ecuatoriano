import { Service, signal } from '@angular/core';

@Service()
export class ArtworkSearchService {
  readonly asset = signal<File | null>(null);
}
