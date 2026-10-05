import { ChangeDetectionStrategy, Component } from '@angular/core';
import { ɵEmptyOutletComponent } from '@angular/router';

@Component({
  selector: 'app-container',
  templateUrl: './container.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class Container {}
