import { ChangeDetectionStrategy, Component, computed, input, output } from '@angular/core';

@Component({
  selector: 'app-outlined-button',
  imports: [],
  templateUrl: './outlined-button.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class OutlinedButton {
  text = input.required<string>();
  color = input<string>('#000000');

  click = output<void>();
}
