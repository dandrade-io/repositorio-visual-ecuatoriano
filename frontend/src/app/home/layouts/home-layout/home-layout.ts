import { Component, computed, effect, inject, input, signal } from '@angular/core';
import { SearchBar } from '../../components/search-bar/search-bar';
import { EmbeddingModel, isEmbeddingModel } from '../../../assets/services/assets.service';
import { AuthService } from '../../../auth/services/auth.service';
import { Router, RouterOutlet } from '@angular/router';
import { Header } from '../../../shared/components/header/header';
import { ArtworkSearchService } from '../../services/artwork-search-service';

@Component({
  selector: 'app-home-layout',
  imports: [SearchBar, Header, RouterOutlet],
  templateUrl: './home-layout.html',
})
export class HomeLayout {
  private router = inject(Router);
  private searchService = inject(ArtworkSearchService);

  protected auth = inject(AuthService);

  query = input<string>('');
  model = input<string>('');

  protected readonly validModel = computed<EmbeddingModel | null>(() => {
    const m = this.model();
    return isEmbeddingModel(m) ? m : null;
  });

  protected onSearch(text: string, asset: File | null, model: EmbeddingModel): void {
    const query = text.trim();
    if (!asset && !query) return;

    this.searchService.asset.set(asset);
    const mode = asset ? 'image' : 'text';
    const q = asset ? null : query;

    this.router.navigate(['/search'], { queryParams: { model, mode: mode, query: q } });
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
