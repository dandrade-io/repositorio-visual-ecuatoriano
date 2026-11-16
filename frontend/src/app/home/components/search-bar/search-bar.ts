import { Component, output, signal } from '@angular/core';
import { EmbeddingModel } from '../../../assets/services/assets.service';
import { KeyValuePipe } from '@angular/common';
import { SearchButton } from '../search-button/search-button';
import { SelectModelButton } from '../select-model-button/select-model-button';

export interface Search {
  query: string;
  model: EmbeddingModel;
}

@Component({
  selector: 'app-search-bar',
  imports: [SearchButton, SelectModelButton],
  templateUrl: './search-bar.html',
})
export class SearchBar {
  protected readonly models: Record<EmbeddingModel, string> = {
    'jina-clip-v2': 'Jina CLIP v2',
    siglip2: 'SigLIP 2',
  };

  readonly search = output<Search>();

  protected readonly text = signal('');
  protected readonly model = signal<EmbeddingModel>('jina-clip-v2');

  protected submit(event: Event) {
    event.preventDefault();

    const query = this.text().trim();
    const model = this.model();
    if (!query) return;

    this.search.emit({ query, model });
  }
}
