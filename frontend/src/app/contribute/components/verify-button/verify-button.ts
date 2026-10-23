import { ChangeDetectionStrategy, Component, computed, input, model } from '@angular/core';
import { VerifierRole } from '../../../auth/services/auth.service';

@Component({
  selector: 'app-verify-button',
  imports: [],
  templateUrl: './verify-button.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class VerifyButton {
  readonly role = input<VerifierRole | null>(null);
  readonly verified = model(false);
  readonly disabled = input(false);

  // Without a verifier role nothing can be verified, regardless of the field.
  protected readonly isDisabled = computed(() => this.disabled() || !this.role());

  protected readonly isExpert = computed(() => this.role() === 'expert_verifier');

  protected readonly classes = computed(() => {
    if (!this.verified()) return 'border-[#d6d1c8] text-stone-400 hover:bg-stone-100';
    return this.isExpert()
      ? 'border-[#17362A]/50 text-[#17362A]'
      : 'border-stone-300 text-stone-700';
  });

  protected toggle() {
    this.verified.update((v) => !v);
  }
}
