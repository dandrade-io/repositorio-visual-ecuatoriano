import { ChangeDetectionStrategy, Component, input } from '@angular/core';
import { ArtworkBase } from '../../interfaces/artwork-base.interface';
import { ArtworkCard } from '../artwork-card/artwork-card';

@Component({
  selector: 'app-artwork-grid',
  imports: [ArtworkCard],
  templateUrl: './artwork-grid.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ArtworkGrid {
  readonly artworks = input.required<ArtworkBase[]>();
}
