import { Component } from '@angular/core';

@Component({
  selector: 'app-search-button',
  imports: [],
  template: `
   <button type="submit" class="cursor-pointer bg-primary h-10 w-10 rounded-full">
    <svg
      class="size-5.25 shrink-0 text-white mx-auto"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
      aria-hidden="true"
    >
      <circle cx="11" cy="11" r="7" />
      <path d="m20 20-3.5-3.5" />
    </svg>
  </button>`,
})
export class SearchButton {}
