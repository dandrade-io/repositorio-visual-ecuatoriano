import {
  ChangeDetectionStrategy,
  Component,
  computed,
  effect,
  inject,
  signal,
  Signal,
} from '@angular/core';
import { AbstractControl, NonNullableFormBuilder, ReactiveFormsModule } from '@angular/forms';
import { InputText } from '../../../contribute/components/input-text/input-text';
import { Container } from '../../../contribute/components/container/container';
import { ChipGroup } from '../../../contribute/components/chip-group/chip-group';
import { Technique, TECHNIQUE_LABELS } from '../../enums/technique.enum';
import { Predominance, PREDOMINANCE_LABELS } from '../../enums/prodominance.enum';
import { TITLE_SYNTAX_LABELS, TitleSyntax } from '../../enums/title-syntax.enum';
import { PostButton } from '../../../contribute/components/post-button/post-button';
import { Button } from '../../../shared/components/button/button';
import { ChipOption } from '../../../contribute/components/chip-item/chip-item';
import { NgClass } from '@angular/common';
import { AlbumCoverService } from '../../services/album-cover.service';
import { rxResource, takeUntilDestroyed, toSignal } from '@angular/core/rxjs-interop';
import { AssetsService } from '../../../assets/services/assets.service';
import { firstValueFrom, map, startWith } from 'rxjs';
import { Router } from '@angular/router';
import { VerifyButton } from '../../../contribute/components/verify-button/verify-button';
import { AuthService } from '../../../auth/services/auth.service';

const VERIFIABLE_FIELDS = [
  'title',
  'year',
  'artist',
  'rhythms',
  'representationTechnique',
  'graphicPredominance',
  'titleSyntax',
] as const;
type VerifiableField = (typeof VERIFIABLE_FIELDS)[number];

@Component({
  selector: 'register-album-cover-page',
  imports: [
    ReactiveFormsModule,
    InputText,
    Container,
    ChipGroup,
    PostButton,
    Button,
    NgClass,
    VerifyButton,
  ],
  templateUrl: './register-album-cover-page.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class RegisterAlbumCoverPage {
  private router = inject(Router);
  private fb = inject(NonNullableFormBuilder);

  protected authService = inject(AuthService);
  private albumCoverService = inject(AlbumCoverService);
  private asssetsService = inject(AssetsService);

  protected form = this.fb.group({
    title: '',
    year: '',
    artist: '',
    rhythms: '',
    representationTechnique: this.fb.control<Technique | null>(null),
    graphicPredominance: this.fb.control<Predominance | null>(null),
    titleSyntax: this.fb.control<TitleSyntax | null>(null),
    // dominantColors: this.fb.control<[]string | null>(null),
  });

  protected readonly verified = signal<Record<VerifiableField, boolean>>(
    Object.fromEntries(VERIFIABLE_FIELDS.map((field) => [field, false])) as Record<
      VerifiableField,
      boolean
    >,
  );

  protected setVerified(field: VerifiableField, value: boolean) {
    this.verified.update((current) => ({ ...current, [field]: value }));
  }

  private readonly formValue = toSignal(
    this.form.valueChanges.pipe(
      map(() => this.form.getRawValue()),
      startWith(this.form.getRawValue()),
    ),
    { requireSync: true },
  );

  // A field can only be verified once it has a value.
  protected readonly hasValue = computed(() => {
    const value = this.formValue();
    return Object.fromEntries(
      VERIFIABLE_FIELDS.map((field) => {
        const v = value[field];
        return [field, typeof v === 'string' ? v.trim() !== '' : v != null];
      }),
    ) as Record<VerifiableField, boolean>;
  });

  protected readonly canVerifyAll = computed(
    () => !!this.authService.verifierRole() && Object.values(this.hasValue()).some(Boolean),
  );

  protected verifyAll() {
    const hasValue = this.hasValue();
    this.verified.update(
      (current) =>
        Object.fromEntries(
          VERIFIABLE_FIELDS.map((field) => [field, current[field] || hasValue[field]]),
        ) as Record<VerifiableField, boolean>,
    );
  }

  constructor() {
    // Editing a field (or autofill / form reset) invalidates its verification.
    for (const field of VERIFIABLE_FIELDS) {
      const control: AbstractControl = this.form.controls[field];
      control.valueChanges
        .pipe(takeUntilDestroyed())
        .subscribe(() => this.setVerified(field, false));
    }

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

  selectedFile = signal<File | null>(null);
  previewUrl = signal<string | null>(null);
  submitAllowed = signal<boolean>(false);

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
  private readonly autofillRequest = signal<{ assetId: string } | null>(null);

  protected readonly albumCoverResource = rxResource({
    params: () => this.autofillRequest() ?? undefined,
    stream: ({ params }) => this.albumCoverService.extractMetadata(params.assetId),
  });

  protected readonly assetResource = rxResource({
    params: () => this.selectedFile() ?? undefined,
    stream: ({ params }) => this.asssetsService.upload(params),
  });

  protected autofillWithAi() {
    const asset = this.assetResource.value();
    if (asset) this.autofillRequest.set({ assetId: asset.id });
  }

  protected async submit() {
    const asset = this.assetResource.value();
    if (!asset) return;

    const { year, rhythms, ...rest } = this.form.getRawValue();
    const payload = {
      ...rest,
      dominantColors: [],
      year: year ? Number(year) : null,
      rhythms: rhythms
        .split(',')
        .map((r) => r.trim())
        .filter(Boolean),
    };

    try {
      this.submitAllowed.set(false);
      await firstValueFrom(this.albumCoverService.upload(payload, asset.id));
      this.router.navigate(['/']);
    } catch (e) {
      console.error(e);
    }

    this.submitAllowed.set(true);
  }

  protected onFileSelected(input: HTMLInputElement) {
    const file: File | undefined = input.files?.[0];
    if (!file) return;

    this.selectedFile.set(file);
    this.submitAllowed.set(true);

    const previous = this.previewUrl();
    if (previous) URL.revokeObjectURL(previous);

    this.previewUrl.set(URL.createObjectURL(file));

    input.value = '';
    this.form.reset();
    this.albumCoverResource.set(undefined);
  }
}
