import { Component, inject } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../../../auth/services/auth.service';
import { Button } from '../button/button';
import { SquarePlus } from '../../icons/square-plus/square-plus';
import { ProfileMenu } from '../profile-menu/profile-menu';

@Component({
  selector: 'app-header',
  host: { class: 'contents' },
  imports: [RouterLink, Button, SquarePlus, ProfileMenu],
  templateUrl: './header.html',
})
export class Header {
  private router = inject(Router);
  protected auth = inject(AuthService);

  register() {
    this.auth.register();
  }

  login() {
    this.auth.login();
  }

  logout() {
    this.auth.logout();
  }

  goToRegisterAlbumCover() {
    this.router.navigate(['/contribute/album-cover']);
  }
}
