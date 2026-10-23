import { ChangeDetectionStrategy, Component, input } from '@angular/core';

@Component({
  selector: 'app-post-button',
  imports: [],
  templateUrl: './post-button.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class PostButton {
  disabled = input<boolean>(false);
}
