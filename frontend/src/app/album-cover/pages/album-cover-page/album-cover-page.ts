import { Component, computed, inject, input } from '@angular/core';
import { rxResource } from '@angular/core/rxjs-interop';
import { AlbumCoverService } from '../../services/album-cover.service';
import { VerifiableField } from '../../enums/verifiable_fields';
import { DotsBounce } from '../../../shared/components/dots-bounce/dots-bounce';
import { AssetsService } from '../../../assets/services/assets.service';
import { TECHNIQUE_LABELS } from '../../enums/technique.enum';
import { PREDOMINANCE_LABELS } from '../../enums/prodominance.enum';
import { TITLE_SYNTAX_LABELS } from '../../enums/title-syntax.enum';
import { Button } from '../../../shared/components/button/button';
import { VerifierRole } from '../../../auth/services/auth.service';
import { VerifierIcon } from '../../../shared/icons/verifier-icon/verifier-icon';

@Component({
  selector: 'app-album-cover-page',
  imports: [DotsBounce, Button, VerifierIcon],
  templateUrl: './album-cover-page.html',
})
export class AlbumCoverPage {
  private readonly albumCoverService = inject(AlbumCoverService);
  private readonly assetsService = inject(AssetsService);

  id = input.required<string>();

  protected readonly albumCoverResource = rxResource({
    params: () => this.id(),
    stream: ({ params: id }) => this.albumCoverService.fetchById(id),
  });

  protected readonly imageUrl = computed(() => {
    if (!this.albumCoverResource.hasValue()) return '';
    const albumCover = this.albumCoverResource.value();

    return this.assetsService.getAssetUrl(albumCover.path);
  });

  protected readonly verifiedBy = computed(() => {
    const checks = this.albumCoverResource.hasValue() ? this.albumCoverResource.value().checks : [];
    const result: Partial<Record<VerifiableField, VerifierRole>> = {};

    for (const { fieldName, verifierRole } of checks) {
      // An expert check always takes priority over an ordinary one.
      if (result[fieldName] !== 'expert_verifier') result[fieldName] = verifierRole;
    }

    return result;
  });

  protected readonly items = computed(() => {
    if (!this.albumCoverResource.hasValue()) return [];
    const albumCover = this.albumCoverResource.value();
    const verifiedBy = this.verifiedBy();

    const rows: { field: VerifiableField; label: string; value: string | number }[] = [
      { field: 'year', label: 'Año', value: albumCover.year ?? 'N.A' },
      {
        field: 'rhythms',
        label: 'Ritmos',
        value: albumCover.rhythms.length != 0 ? albumCover.rhythms.join(', ') : 'N.A',
      },
      {
        field: 'representationTechnique',
        label: 'Técnica de Representación',
        value: albumCover.representationTechnique
          ? TECHNIQUE_LABELS[albumCover.representationTechnique]
          : 'N.A',
      },
      {
        field: 'predominance',
        label: 'Predominancia gráfica',
        value: albumCover.predominance ? PREDOMINANCE_LABELS[albumCover.predominance] : 'N.A',
      },
      {
        field: 'titleSyntax',
        label: 'Sintaxis del título',
        value: albumCover.titleSyntax ? TITLE_SYNTAX_LABELS[albumCover.titleSyntax] : 'N.A',
      },
    ];

    return rows.map((row) => ({ ...row, verifiedBy: verifiedBy[row.field] ?? null }));
  });

  protected isExpert(role: VerifierRole | undefined) {
    return role === 'expert_verifier';
  }

  protected retryFetch() {
    this.albumCoverResource.reload();
  }
}
