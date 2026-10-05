import { ChangeDetectionStrategy, Component, computed, input, output } from '@angular/core';

export interface ChipOption<T> {
  value: T;
  label: string;
}

@Component({
  selector: 'app-chip-item',
  imports: [],
  templateUrl: './chip-item.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ChipItem<T> {
  option = input.required<ChipOption<T>>();
  isSelected = input<boolean>(false);

  click = output<void>();

  stateClasses = computed(() =>
    this.isSelected()
      ? 'border-[#17362A] bg-[#17362A] text-white'
      : 'border-[#E7E1D8] text-[#1C1917] hover:bg-[#E7E1D8]/30',
  );

  onPressed() {
    this.click.emit();
  }
}
