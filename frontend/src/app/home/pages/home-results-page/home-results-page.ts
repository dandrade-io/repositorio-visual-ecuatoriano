import { Component, computed, inject, input } from '@angular/core';
import { rxResource } from '@angular/core/rxjs-interop';
import { AlbumCoverService } from '../../../album-cover/services/album-cover.service';
import {
  AssetsService,
  EmbeddingModel,
  isEmbeddingModel,
} from '../../../assets/services/assets.service';
import { map } from 'rxjs';
import { ArtworkMapper } from '../../../artwork/components/mappers/artwork.mapper';
import { ArtworkGrid } from '../../../artwork/components/artwork-grid/artwork-grid';
import { ArtworkSearchService } from '../../services/artwork-search-service';

@Component({
  selector: 'app-home-results-page',
  imports: [ArtworkGrid],
  templateUrl: './home-results-page.html',
})
export class HomeResultsPage {
  private assetsService = inject(AssetsService);
  private albumCoverService = inject(AlbumCoverService);
  private searchService = inject(ArtworkSearchService);

  query = input<string>('');
  model = input<string>('');

  protected readonly validModel = computed<EmbeddingModel | null>(() => {
    const m = this.model();
    return isEmbeddingModel(m) ? m : null;
  });

  protected readonly artworksResource = rxResource({
    params: () => ({
      query: this.query(),
      model: this.validModel(),
      asset: this.searchService.asset(),
    }),
    stream: ({ params }) => {
      if ((!params.query && !params.asset) || !params.model) {
        return this.albumCoverService.fetchSummary(10);
      }

      console.log(params.asset);

      if (params.asset) {
        return this.assetsService
          .searchByAsset(params.asset, 'image', params.model, 10)
          .pipe(map(ArtworkMapper.fromAssetSearchResponse));
      }

      return this.assetsService
        .search({ q: params.query, kind: 'image', model: params.model, limit: 10 })
        .pipe(map(ArtworkMapper.fromAssetSearchResponse));
    },
  });
}
