import { Component } from '@angular/core';

@Component({
  selector: 'app-dots-bounce',
  imports: [],
  template: `
    <div class="flex gap-2" role="status" aria-label="Loading">
      <span
        class="size-3.5 animate-dot-wave rounded-full bg-[#155dfc] [animation-delay:-0.3s] motion-reduce:animate-none"
        aria-hidden="true"
      ></span>
      <span
        class="size-3.5 animate-dot-wave rounded-full bg-[#155dfc] [animation-delay:-0.15s] motion-reduce:animate-none"
        aria-hidden="true"
      ></span>
      <span
        class="size-3.5 animate-dot-wave rounded-full bg-[#155dfc] motion-reduce:animate-none"
        aria-hidden="true"
      ></span>
    </div>
  `,
})
export class DotsBounce {}
