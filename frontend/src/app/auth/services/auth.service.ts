import { computed, effect, inject, Injectable, signal } from '@angular/core';
import {
  KEYCLOAK_EVENT_SIGNAL,
  KeycloakEventType,
  ReadyArgs,
  typeEventArgs,
} from 'keycloak-angular';
import Keycloak from 'keycloak-js';

export type VerifierRole = 'expert_verifier' | 'ordinary_verifier';
const VERIFIER_ROLES: VerifierRole[] = ['expert_verifier', 'ordinary_verifier'];

@Injectable({ providedIn: 'root' })
export class AuthService {
  private keycloak = inject(Keycloak);
  private events = inject(KEYCLOAK_EVENT_SIGNAL);

  private _authenticated = signal(this.keycloak.authenticated ?? false);
  readonly authenticated = this._authenticated.asReadonly();

  constructor() {
    effect(() => {
      const event = this.events();

      switch (event.type) {
        case KeycloakEventType.Ready:
          this._authenticated.set(typeEventArgs<ReadyArgs>(event.args));
          break;
        case KeycloakEventType.AuthLogout:
        case KeycloakEventType.AuthRefreshError:
        case KeycloakEventType.TokenExpired:
          this._authenticated.set(this.keycloak.authenticated ?? false);
          break;
      }
    });
  }

  readonly verifierRole = computed<VerifierRole | null>(() => {
    this.authenticated();

    const roles: string[] =
      this.keycloak.tokenParsed?.resource_access?.['backend-api']?.roles ?? [];

    return VERIFIER_ROLES.find((role) => roles.includes(role)) ?? null;
  });

  login() {
    this.keycloak.login();
  }

  register() {
    this.keycloak.register();
  }

  async logout() {
    await this.keycloak.logout();
  }
}
