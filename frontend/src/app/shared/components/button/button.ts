import { ChangeDetectionStrategy, Component, computed, input } from '@angular/core';

export type ButtonVariant = 'solid' | 'outlined';

@Component({
  selector: 'app-button',
  imports: [],
  templateUrl: './button.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class Button {
  type = input<string>('button');
  disabled = input<boolean>(false);
  variant = input<ButtonVariant>('solid');
  color = input<string>('#155dfc');

  protected readonly isOutlined = computed(() => this.variant() === 'outlined');

  safeColor = computed(() => (/^#[0-9a-fA-F]{6}$/.test(this.color()) ? this.color() : '#64748b'));
  textColor = computed(() => {
    const h = this.safeColor().slice(1);
    const [r, g, b] = [0, 2, 4].map((i) => parseInt(h.slice(i, i + 2), 16));
    return (r * 299 + g * 587 + b * 114) / 1000 > 150 ? '#111827' : '#ffffff';
  });
}
