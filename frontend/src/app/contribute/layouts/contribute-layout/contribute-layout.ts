import { ChangeDetectionStrategy, Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';

@Component({
  selector: 'contribute-layout',
  imports: [RouterOutlet],
  templateUrl: './contribute-layout.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ContributeLayout {}
