import { Component } from '@angular/core';

@Component({
  selector: 'app-search-button',
  imports: [],
  template: `
    <button
      type="submit"
      class="cursor-pointer bg-primary transition hover:bg-neutral-800 h-10 w-10 rounded-full pr-px"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 20 20"
        fill="currentColor"
        class="size-5 text-white mx-auto mx"
      >
        <path
          fill-rule="evenodd"
          d="M9 3.5a5.5 5.5 0 1 0 0 11 5.5 5.5 0 0 0 0-11ZM2 9a7 7 0 1 1 12.452 4.391l3.328 3.329a.75.75 0 1 1-1.06 1.06l-3.329-3.328A7 7 0 0 1 2 9Z"
          clip-rule="evenodd"
1        />
      </svg>
    </button>
  `,
})
export class SearchButton {}
