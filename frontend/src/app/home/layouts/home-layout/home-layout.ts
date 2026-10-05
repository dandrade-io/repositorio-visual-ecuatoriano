import { ChangeDetectionStrategy, Component, inject, signal } from '@angular/core';
import { SearchBar } from '../../components/search-bar/search-bar';
import { AssetsService } from '../../../assets/services/assets.service';
import { rxResource } from '@angular/core/rxjs-interop';

@Component({
  selector: 'app-home-layout',
  imports: [SearchBar],
  templateUrl: './home-layout.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class HomeLayout {
  private assetsService = inject(AssetsService);
  protected readonly imageUrl = (id: string | number) => this.assetsService.imageUrl(id);
  protected readonly query = signal<string | undefined>(undefined);

  protected readonly assetsResource = rxResource({
    params: () => this.query(),
    stream: ({ params }) => this.assetsService.search({q: params, kind: "image", model: "jina-clip-v2", limit: 10}),
  });

  protected onSearch(text: string): void {
    this.query.set(text);
  }
}
