import { ChangeDetectionStrategy, Component, computed, inject, input } from '@angular/core';
import { ArtworkBase } from '../../interfaces/artwork-base.interface';
import { EnvironmentConfigService } from '../../../shared/services/environment-config.service';
import { DecimalPipe } from '@angular/common';
import { Router } from '@angular/router';

@Component({
  selector: 'app-artwork-card',
  imports: [DecimalPipe],
  templateUrl: './artwork-card.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ArtworkCard {
  private router = inject(Router);

  protected readonly env = inject(EnvironmentConfigService);

  readonly artwork = input.required<ArtworkBase>();

  protected readonly imageUrl = computed(() => {
    const apiUrl = this.env.get(EnvironmentConfigService.API_URL);
    const artwork = this.artwork();

    return `${apiUrl}/assets/${artwork.imageUrl}`;
  });

  protected navigateToDetailsPage() {
    const artwork = this.artwork();

    if (artwork.type === 'album-cover') {
      this.router.navigate(['/artwork', artwork.type, artwork.id]);
    }
  }
}
