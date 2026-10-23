import { CdkMenu, CdkMenuItem, CdkMenuTrigger } from '@angular/cdk/menu';
import { ChangeDetectionStrategy, Component, computed, inject } from '@angular/core';
import { AuthService, VerifierRole } from '../../../auth/services/auth.service';
import { RouterLink } from '@angular/router';
import Keycloak from 'keycloak-js';

export interface RoleBadge {
  label: string;
  classes: string;
}

export interface ProfileInformation {
  name: string;
  username: string;
  role: RoleBadge | null;
}

const ROLE_BADGES: Record<VerifierRole, RoleBadge> = {
  expert_verifier: {
    label: 'Verificador Experto',
    classes: 'bg-[#17362A]/10 text-[#17362A]',
  },
  ordinary_verifier: {
    label: 'Verificador Ordinario',
    classes: 'bg-stone-100 text-stone-700',
  },
};

@Component({
  selector: 'app-profile-menu',
  imports: [CdkMenuTrigger, CdkMenu, CdkMenuItem],
  templateUrl: './profile-menu.html',
})
export class ProfileMenu {
  private authService = inject(AuthService);
  private keycloak = inject(Keycloak);

  protected profileInfo = computed<ProfileInformation>(() => {
    const token = this.keycloak.tokenParsed;
    const roleKey = this.authService.verifierRole();

    return {
      name: token?.['name'] ?? '',
      username: token?.['preferred_username'] ?? '',
      role: roleKey ? ROLE_BADGES[roleKey] : null,
    };
  });

  protected logout() {
    this.authService.logout();
  }
}
