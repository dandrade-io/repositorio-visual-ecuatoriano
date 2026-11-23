import { Component, model, output } from '@angular/core';
import { EmbeddingModel } from '../../../assets/services/assets.service';
import { KeyValuePipe } from '@angular/common';

@Component({
  selector: 'app-select-model-button',
  imports: [KeyValuePipe],
  templateUrl: './select-model-button.html',
})
export class SelectModelButton {
  protected readonly models: Record<EmbeddingModel, string> = {
    'jina-clip-v2': 'Jina CLIP v2',
    siglip2: 'SigLIP 2',
    'clamp3': 'CLaMP3'
  };

  readonly selected = model.required<EmbeddingModel>();

  protected onModelSelected(model: string) {
    this.selected.set(model as EmbeddingModel);
  }
}
