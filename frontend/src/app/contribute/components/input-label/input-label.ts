import { ChangeDetectionStrategy, Component, input } from '@angular/core';

@Component({
  selector: 'input-label',
  imports: [],
  templateUrl: './input-label.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class InputLabel {
  for = input.required<string>();
  label = input.required<string>();
}
