import { Component, computed, input } from '@angular/core';

@Component({
  selector: 'app-verifier-icon',
  imports: [],
  host: { class: 'inline-flex' },
  template: `
    <span
      class="group relative inline-flex"
      [attr.tabindex]="showTooltip() ? 0 : null"
      [attr.aria-label]="showTooltip() ? tooltipText() : null"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        fill="currentColor"
        class="size-5"
        [class.text-[#17362A]]="isExpert()"
        [class.text-stone-700]="!isExpert()"
        aria-hidden="true"
      >
        @if (isExpert()) {
          <path
            fill-rule="evenodd"
            d="M12.516 2.17a.75.75 0 0 0-1.032 0 11.209 11.209 0 0 1-7.877 3.08.75.75 0 0 0-.722.515A12.74 12.74 0 0 0 2.25 9.75c0 5.942 4.064 10.933 9.563 12.348a.749.749 0 0 0 .374 0c5.499-1.415 9.563-6.406 9.563-12.348 0-1.39-.223-2.73-.635-3.985a.75.75 0 0 0-.722-.516l-.143.001c-2.996 0-5.717-1.17-7.734-3.08Zm3.094 8.016a.75.75 0 1 0-1.22-.872l-3.236 4.53L9.53 12.22a.75.75 0 0 0-1.06 1.06l2.25 2.25a.75.75 0 0 0 1.14-.094l3.75-5.25Z"
            clip-rule="evenodd"
          />
        } @else {
          <path
            fill-rule="evenodd"
            d="M2.25 12c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75-4.365 9.75-9.75 9.75S2.25 17.385 2.25 12Zm13.36-1.814a.75.75 0 1 0-1.22-.872l-3.236 4.53L9.53 12.22a.75.75 0 0 0-1.06 1.06l2.25 2.25a.75.75 0 0 0 1.14-.094l3.75-5.25Z"
            clip-rule="evenodd"
          />
        }
      </svg>
      @if (showTooltip()) {
        <span
          role="tooltip"
          class="pointer-events-none absolute bottom-full left-1/2 z-10 mb-1 -translate-x-1/2 whitespace-nowrap rounded bg-gray-900 px-2 py-1 text-xs font-normal text-white opacity-0 transition-opacity group-hover:opacity-100 group-focus:opacity-100"
        >
          {{ tooltipText() }}
        </span>
      }
    </span>
  `,
})
export class VerifierIcon {
  readonly isExpert = input(false);
  readonly showTooltip = input(false);

  protected readonly tooltipText = computed(() =>
    this.isExpert() ? 'Verificado por un usuario experto' : 'Verificado por un usuario ordinario',
  );
}
