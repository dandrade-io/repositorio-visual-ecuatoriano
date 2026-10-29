import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { Header } from '../../../shared/components/header/header';

@Component({
  selector: 'app-artwork-layout',
  imports: [RouterOutlet, Header],
  template: `
    <div class="flex flex-col min-h-dvh bg-[#F8F8FD]">
      <app-header />
      <main class="flex flex-1 justify-center px-[12vw] py-8 bg-[#F8F8FD]">
        <router-outlet />
      </main>
    </div>
  `,
})
export class ArtworkLayout {}
