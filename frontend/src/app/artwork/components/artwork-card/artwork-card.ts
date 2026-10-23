import { ChangeDetectionStrategy, Component, computed, inject, input } from '@angular/core';
import { ArtworkBase } from '../../interfaces/artwork-base.interface';
import { EnvironmentConfigService } from '../../../shared/services/environment-config.service';
import { DecimalPipe } from '@angular/common';

@Component({
  selector: 'app-artwork-card',
  imports: [DecimalPipe],
  templateUrl: './artwork-card.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ArtworkCard {
  protected readonly env = inject(EnvironmentConfigService);

  readonly artwork = input.required<ArtworkBase>();

  protected readonly imageUrl = computed(() => {
    const apiUrl = this.env.get(EnvironmentConfigService.API_URL);
    const artwork = this.artwork();

    return `${apiUrl}/assets/${artwork.imageUrl}`;
  });
}
