import { Component, model, output, signal } from '@angular/core';
import { EmbeddingModel } from '../../../assets/services/assets.service';
import { SearchButton } from '../search-button/search-button';
import { SelectModelButton } from '../select-model-button/select-model-button';
import { AttachmentButton } from '../attachment-button/attachment-button';
import { AttachmentPreview } from '../attachment-preview/attachment-preview';

export interface Search {
  query: string;
  asset?: File | null;
  model: EmbeddingModel;
}

@Component({
  selector: 'app-search-bar',
  imports: [SearchButton, SelectModelButton, AttachmentButton, AttachmentPreview],
  templateUrl: './search-bar.html',
})
export class SearchBar {
  protected readonly models: Record<EmbeddingModel, string> = {
    'jina-clip-v2': 'Jina CLIP v2',
    siglip2: 'SigLIP 2',
    'clamp3': 'CLaMP3'
  };

  readonly asset = signal<File | null>(null);

  readonly search = output<Search>();

  readonly query = model<string>('');
  readonly model = model.required<EmbeddingModel>();

  protected submit(event: Event) {
    event.preventDefault();

    const query = this.query().trim();
    const asset = this.asset();
    const model = this.model();
    if (!query && !asset) return;

    this.search.emit({ query, asset, model });
  }
}
