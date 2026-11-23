import { ChangeDetectionStrategy, Component, computed, inject, input, signal } from '@angular/core';
import { ArtworkBase } from '../../interfaces/artwork-base.interface';
import { EnvironmentConfigService } from '../../../shared/services/environment-config.service';
import { DecimalPipe } from '@angular/common';
import { Router } from '@angular/router';
import { AudioPlayer } from '../audio-player/audio-player';

@Component({
  selector: 'app-artwork-card',
  imports: [DecimalPipe, AudioPlayer],
  templateUrl: './artwork-card.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ArtworkCard {
  private router = inject(Router);

  protected readonly env = inject(EnvironmentConfigService);

  readonly artwork = input.required<ArtworkBase>();

  protected readonly imgLoaded = signal(false);
  protected readonly imgFailed = signal(false);

  protected readonly mediaUrl = computed(() => {
    const apiUrl = this.env.get(EnvironmentConfigService.API_URL);
    const artwork = this.artwork();

    return `${apiUrl}/assets/${artwork.url}`;
  });

  protected readonly isAudio = computed(
    () => this.artwork().type === 'assets' && this.artwork().kind === 'audio',
  );

  protected navigateToDetailsPage() {
    const artwork = this.artwork();

    if (artwork.type === 'album-cover') {
      this.router.navigate(['/artwork', artwork.type, artwork.id]);
    }
  }
}
