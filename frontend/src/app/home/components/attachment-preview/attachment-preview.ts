import { Component, input, output, signal } from '@angular/core';

@Component({
  selector: 'app-attachment-preview',
  imports: [],
  template: `
    <div class="flex items-center bg-[#E5E7EB] px-3 rounded-lg py-1.5">
      <span class="block max-w-48 truncate text-sm">{{ name() }}</span>
      <div class="w-1.5"></div>

      <button (click)="remove.emit()" class="cursor-pointer">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 20 20"
          fill="currentColor"
          class="size-4.25 text-primary hover:text-neutral-700"
        >
          <path
            fill-rule="evenodd"
            d="M10 18a8 8 0 1 0 0-16 8 8 0 0 0 0 16ZM8.28 7.22a.75.75 0 0 0-1.06 1.06L8.94 10l-1.72 1.72a.75.75 0 1 0 1.06 1.06L10 11.06l1.72 1.72a.75.75 0 1 0 1.06-1.06L11.06 10l1.72-1.72a.75.75 0 0 0-1.06-1.06L10 8.94 8.28 7.22Z"
            clip-rule="evenodd"
          />
        </svg>
      </button>
    </div>
  `,
})
export class AttachmentPreview {
  readonly name = input<string>();

  readonly remove = output<void>();
}
