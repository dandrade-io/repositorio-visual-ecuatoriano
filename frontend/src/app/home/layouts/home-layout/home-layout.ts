import { Component, inject, signal } from '@angular/core';
import { SearchBar } from '../../components/search-bar/search-bar';
import { AssetsService } from '../../../assets/services/assets.service';
import { rxResource } from '@angular/core/rxjs-interop';
import { AuthService } from '../../../auth/services/auth.service';
import { Router } from '@angular/router';
import { map } from 'rxjs';
import { ArtworkMapper } from '../../../artwork/components/mappers/artwork.mapper';
import { ArtworkGrid } from '../../../artwork/components/artwork-grid/artwork-grid';
import { AlbumCoverService } from '../../../album-cover/services/album-cover.service';
import { Header } from '../../../shared/components/header/header';

@Component({
  selector: 'app-home-layout',
  imports: [SearchBar, ArtworkGrid, Header],
  templateUrl: './home-layout.html',
})
export class HomeLayout {
  private router = inject(Router);

  private albumCoverService = inject(AlbumCoverService);
  private assetsService = inject(AssetsService);
  protected auth = inject(AuthService);

  protected readonly query = signal<string>('');

  protected readonly artworksResource = rxResource({
    params: () => this.query(),
    stream: ({ params }) => {
      // if no query value, fetch any album cover
      if (!params.trim()) {
        return this.albumCoverService
          .fetchSummary()
          .pipe(map(ArtworkMapper.fromAlbumCoverSummaryResponse));
      }

      return this.assetsService
        .search({ q: params, kind: 'image', model: 'jina-clip-v2', limit: 10 })
        .pipe(map(ArtworkMapper.fromAssetSearchResponse));
    },
  });

  onSearch(text: string): void {
    this.query.set(text);
  }

  register() {
    this.auth.register();
  }

  login() {
    this.auth.login();
  }

  logout() {
    this.auth.logout();
  }

  goToRegisterAlbumCoverPage() {
    this.router.navigate(['/contribute/album-cover']);
  }
}
