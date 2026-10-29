import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';

@Component({
  selector: 'contribute-layout',
  imports: [RouterOutlet],
  template: `
    <div
      class="flex justify-center min-h-screen min-w-screen px-[15vw] lg:px-[20vw] py-9 bg-[#F8F8FD]"
    >
      <router-outlet />
    </div>
  `,
})
export class ContributeLayout {}
