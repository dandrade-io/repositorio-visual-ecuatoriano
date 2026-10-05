import { ChangeDetectionStrategy, Component, forwardRef, input, signal } from '@angular/core';
import { ControlValueAccessor, NG_VALUE_ACCESSOR } from '@angular/forms';
import { ChipItem, ChipOption } from '../chip-item/chip-item';
import { InputLabel } from '../input-label/input-label';

@Component({
  selector: 'app-chip-group',
  imports: [ChipItem, InputLabel],
  templateUrl: './chip-group.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
  providers: [
    { provide: NG_VALUE_ACCESSOR, useExisting: forwardRef(() => ChipGroup), multi: true },
  ],
})
export class ChipGroup<T> implements ControlValueAccessor {
  label = input<string>('');
  options = input.required<ChipOption<T>[]>();

  protected value = signal<T | null>(null);
  protected disabled = signal(false);

  private onChange: (value: T | null) => void = () => {};
  private onTouched: () => void = () => {};

  writeValue(value: T | null): void {
    this.value.set(value ?? null);
  }

  registerOnChange(fn: (value: T | null) => void): void {
    this.onChange = fn;
  }

  registerOnTouched(fn: () => void): void {
    this.onTouched = fn;
  }

  setDisabledState(isDisabled: boolean): void {
    this.disabled.set(isDisabled);
  }

  onChipItemPressed(value: T) {
    if (this.disabled()) return;
    this.value.set(value);
    this.onChange(value);
    this.onTouched();
  }
}
