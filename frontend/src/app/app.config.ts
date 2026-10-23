import {
  ApplicationConfig,
  EnvironmentInjector,
  inject,
  provideAppInitializer,
  provideBrowserGlobalErrorListeners,
} from '@angular/core';
import { provideRouter } from '@angular/router';

import { routes } from './app.routes';
import { provideHttpClient, withFetch, withInterceptors } from '@angular/common/http';
import { EnvironmentConfigService } from './shared/services/environment-config.service';
import Keycloak from 'keycloak-js';
import { firstValueFrom } from 'rxjs';
import { createKeycloakSignal, KEYCLOAK_EVENT_SIGNAL } from 'keycloak-angular';
import { bearerInterceptor } from './auth/interceptors/bearer.interceptor';

export const appConfig: ApplicationConfig = {
  providers: [
    provideBrowserGlobalErrorListeners(),
    provideRouter(routes),
    provideHttpClient(withFetch(), withInterceptors([bearerInterceptor])),
    {
      provide: Keycloak,
      useFactory: () => {
        const env = inject(EnvironmentConfigService);
        const url = env.get(EnvironmentConfigService.KC_URL);
        if (!url) throw new Error('Keycloak created before env config was loaded');

        return new Keycloak({
          url: url,
          realm: env.get(EnvironmentConfigService.KC_REALM),
          clientId: env.get(EnvironmentConfigService.KC_CLIENT_ID),
        });
      },
    },
    {
      provide: KEYCLOAK_EVENT_SIGNAL,
      useFactory: () => createKeycloakSignal(inject(Keycloak)),
    },
    provideAppInitializer(async () => {
      const injector = inject(EnvironmentInjector);
      await firstValueFrom(injector.get(EnvironmentConfigService).load());

      const keycloak = injector.get(Keycloak);
      injector.get(KEYCLOAK_EVENT_SIGNAL);

      await keycloak
        .init({
          onLoad: 'check-sso',
          silentCheckSsoRedirectUri: window.location.origin + '/silent-check-sso.html',
        })
        .catch((err) => console.error('Keycloak initialization failed', err));
    }),
  ],
};
