import { Component, computed, inject, input } from '@angular/core';
import { rxResource } from '@angular/core/rxjs-interop';
import { AlbumCoverService } from '../../services/album-cover.service';
import { DotsBounce } from '../../../shared/components/dots-bounce/dots-bounce';
import { AssetsService } from '../../../assets/services/assets.service';
import { TECHNIQUE_LABELS } from '../../enums/technique.enum';
import { PREDOMINANCE_LABELS } from '../../enums/prodominance.enum';
import { TITLE_SYNTAX_LABELS } from '../../enums/title-syntax.enum';
import { Button } from '../../../shared/components/button/button';

@Component({
  selector: 'app-album-cover-page',
  imports: [DotsBounce, Button],
  templateUrl: './album-cover-page.html',
})
export class AlbumCoverPage {
  private readonly albumCoverService = inject(AlbumCoverService);
  private readonly assetsService = inject(AssetsService);

  id = input.required<string>();

  protected readonly imageUrl = computed(() => {
    if (!this.albumCoverResource.hasValue()) return '';
    const albumCover = this.albumCoverResource.value();

    return this.assetsService.getAssetUrl(albumCover.path);
  });

  protected readonly items = computed(() => {
    if (!this.albumCoverResource.hasValue()) return [];
    const albumCover = this.albumCoverResource.value();

    const values = [
      albumCover.year ?? 'N.A',
      albumCover.rhythms.length != 0 ? albumCover.rhythms.join(', ') : 'N.A',
      albumCover.representationTechnique
        ? TECHNIQUE_LABELS[albumCover.representationTechnique]
        : 'N.A',
      albumCover.graphicPredominance ? PREDOMINANCE_LABELS[albumCover.graphicPredominance] : 'N.A',
      albumCover.titleSyntax ? TITLE_SYNTAX_LABELS[albumCover.titleSyntax] : 'N.A',
    ];

    return this.labels.map((label, i) => ({ label, value: values[i] }));
  });

  protected readonly albumCoverResource = rxResource({
    params: () => this.id(),
    stream: ({ params: id }) => this.albumCoverService.fetchById(id),
  });

  private readonly labels = [
    'Año',
    'Ritmos',
    'Técnica de Representación',
    'Predominancia gráfica',
    'Sintaxis del título',
  ];

  protected retryFetch() {
    this.albumCoverResource.reload();
  }
}
