import { Injectable, signal } from '@angular/core';

/** Ensures only one audio player plays at a time. */
@Injectable({ providedIn: 'root' })
export class AudioPlaybackService {
  private readonly _activePlayerId = signal<symbol | null>(null);

  readonly activePlayerId = this._activePlayerId.asReadonly();

  activate(id: symbol) {
    this._activePlayerId.set(id);
  }

  release(id: symbol) {
    if (this._activePlayerId() === id) {
      this._activePlayerId.set(null);
    }
  }
}
