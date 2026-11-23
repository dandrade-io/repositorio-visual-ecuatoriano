import { Component, output, signal } from '@angular/core';

@Component({
  selector: 'app-attachment-button',
  imports: [],
  template: `
    <input
      #inputFile
      type="file"
      class="hidden"
      accept="image/*"
      (change)="onFileSelected(inputFile)"
    />
    <button
      type="button"
      class="bg-[#F2735B] transition hover:brightness-115 h-10 w-10 rounded-full"
      (click)="inputFile.click()"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 20 20"
        fill="currentColor"
        class="size-5 text-white mx-auto"
      >
        <path
          fill-rule="evenodd"
          d="M15.621 4.379a3 3 0 0 0-4.242 0l-7 7a3 3 0 0 0 4.241 4.243h.001l.497-.5a.75.75 0 0 1 1.064 1.057l-.498.501-.002.002a4.5 4.5 0 0 1-6.364-6.364l7-7a4.5 4.5 0 0 1 6.368 6.36l-3.455 3.553A2.625 2.625 0 1 1 9.52 9.52l3.45-3.451a.75.75 0 1 1 1.061 1.06l-3.45 3.451a1.125 1.125 0 0 0 1.587 1.595l3.454-3.553a3 3 0 0 0 0-4.242Z"
          clip-rule="evenodd"
        />
      </svg>
    </button>
  `,
})
export class AttachmentButton {
  // private readonly previousUrl = signal<string | null>(null);

  readonly selected = output<File>();

  private onFileSelected(input: HTMLInputElement) {
    const file: File | undefined = input.files?.[0];
    if (!file) return;

    this.selected.emit(file);

    // const previous = this.previousUrl();
    // if (previous) URL.revokeObjectURL(previous);
    // const url = URL.createObjectURL(file);
  }
}
