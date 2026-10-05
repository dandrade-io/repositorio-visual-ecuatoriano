import {
  ChangeDetectionStrategy,
  Component,
  computed,
  effect,
  inject,
  signal,
  Signal,
} from '@angular/core';
import { NonNullableFormBuilder, ReactiveFormsModule } from '@angular/forms';
import { InputText } from '../../../contribute/components/input-text/input-text';
import { Container } from '../../../contribute/components/container/container';
import { ChipGroup } from '../../../contribute/components/chip-group/chip-group';
import { Technique, TECHNIQUE_LABELS } from '../../enums/technique.enum';
import { Predominance, PREDOMINANCE_LABELS } from '../../enums/prodominance.enum';
import { TITLE_SYNTAX_LABELS, TitleSyntax } from '../../enums/title-syntax.enum';
import { PostButton } from '../../../contribute/components/post-button/post-button';
import { ChipOption } from '../../../contribute/components/chip-item/chip-item';
import { OutlinedButton } from '../../../shared/components/outlined-button/outlined-button';
import { NgClass } from '@angular/common';
import { AlbumCoverService } from '../../services/album-cover.service';
import { rxResource } from '@angular/core/rxjs-interop';
import { tap } from 'rxjs';

@Component({
  selector: 'register-album-cover-page',
  imports: [
    ReactiveFormsModule,
    InputText,
    Container,
    ChipGroup,
    PostButton,
    OutlinedButton,
    NgClass,
  ],
  templateUrl: './register-album-cover-page.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class RegisterAlbumCoverPage {
  private readonly albumCoverService = inject(AlbumCoverService);
  private readonly fb = inject(NonNullableFormBuilder);

  protected readonly form = this.fb.group({
    title: '',
    year: '',
    artist: '',
    rhythms: '',
    representationTechnique: this.fb.control<Technique | null>(null),
    graphicPredominance: this.fb.control<Predominance | null>(null),
    titleSyntax: this.fb.control<TitleSyntax | null>(null),
  });

  selectedFile = signal<File | null>(null);
  previewUrl = signal<string | null>(null);

  techniqueOptions: Signal<ChipOption<Technique>[]> = computed(() =>
    Object.values(Technique).map((value) => ({ value, label: TECHNIQUE_LABELS[value] })),
  );

  predominanceOptions: Signal<ChipOption<Predominance>[]> = computed(() =>
    Object.values(Predominance).map((value) => ({ value, label: PREDOMINANCE_LABELS[value] })),
  );

  titleSyntaxOptions: Signal<ChipOption<TitleSyntax>[]> = computed(() =>
    Object.values(TitleSyntax).map((value) => ({ value, label: TITLE_SYNTAX_LABELS[value] })),
  );

  // Wrapped in an object so that every click triggers a new request, even for the same file.
  private readonly autofillRequest = signal<{ file: File } | null>(null);

  protected readonly albumCoverResource = rxResource({
    params: () => this.autofillRequest() ?? undefined,
    stream: ({ params }) => this.albumCoverService.extractMetadata(params.file),
  });

  constructor() {
    effect(() => {
      const cover = this.albumCoverResource.value();
      if (!cover) return;

      this.form.patchValue({
        title: cover.title ?? '',
        year: cover.year?.toString() ?? '',
        artist: cover.artist ?? '',
        rhythms: cover.rhythms.join(', '),
        representationTechnique: cover.representationTechnique,
        graphicPredominance: cover.graphicPredominance,
        titleSyntax: cover.titleSyntax,
      });
    });
  }

  protected autofillWithAi() {
    const file = this.selectedFile();
    if (file) this.autofillRequest.set({ file });
  }

  protected submit() {
    const { year, rhythms, ...rest } = this.form.getRawValue();
    const payload = {
      ...rest,
      year: year ? Number(year) : null,
      rhythms: rhythms
        .split(',')
        .map((r) => r.trim())
        .filter(Boolean),
    };
    // TODO: send `payload` and `selectedFile()` to the backend once the endpoint exists.
    console.log(payload);
  }

  protected onFileSelected(input: HTMLInputElement) {
    const file: File | undefined = input.files?.[0];
    if (!file) return;

    this.selectedFile.set(file);

    const previous = this.previewUrl();
    if (previous) URL.revokeObjectURL(previous);
    this.previewUrl.set(URL.createObjectURL(file));

    input.value = '';
    this.form.reset();
    this.albumCoverResource.set(undefined);
  }
}
