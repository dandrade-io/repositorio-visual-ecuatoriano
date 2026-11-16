import { Component, computed, input } from '@angular/core';

export type ButtonVariant = 'solid' | 'outlined';

@Component({
  selector: 'app-button',
  imports: [],
  templateUrl: './button.html',
})
export class Button {
  type = input<string>('button');
  disabled = input<boolean>(false);
  variant = input<ButtonVariant>('solid');
  color = input<string | undefined>(undefined);
  tone = input<'primary' | 'custom'>('primary');

  protected readonly isOutlined = computed(() => this.variant() === 'outlined');

  private readonly hex = computed(() => {
    const c = this.color();
    return c && /^#[0-9a-fA-F]{6}$/.test(c) ? c : '#64748b';
  });

  private readonly hexText = computed(() => {
    const h = this.hex().slice(1);
    const [r, g, b] = [0, 2, 4].map((i) => parseInt(h.slice(i, i + 2), 16));
    return (r * 299 + g * 587 + b * 114) / 1000 > 150 ? '#111827' : '#ffffff';
  });

  protected readonly btnColor = computed(() =>
    this.color() === undefined ? 'var(--primary)' : this.hex(),
  );

  protected readonly btnText = computed(() =>
    this.color() === undefined ? 'var(--primary-foreground)' : this.hexText(),
  );
}
