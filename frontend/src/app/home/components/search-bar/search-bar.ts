import { ChangeDetectionStrategy, Component, output, signal } from '@angular/core';
import { Button } from '../../../shared/components/button/button';

@Component({
  selector: 'search-bar',
  imports: [Button],
  templateUrl: './search-bar.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SearchBar {
  readonly search = output<string>();

  protected readonly text = signal('');

  protected submit(event: Event) {
    event.preventDefault();

    const query = this.text().trim();
    if (!query) return;

    this.search.emit(query);
  }
}
