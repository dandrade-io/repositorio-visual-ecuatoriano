import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import Keycloak from 'keycloak-js';
import { EnvironmentConfigService } from '../../shared/services/environment-config.service';
import { from, switchMap } from 'rxjs';

export const bearerInterceptor: HttpInterceptorFn = (req, next) => {
  const env = inject(EnvironmentConfigService);

  const apiUrl = env.get(EnvironmentConfigService.API_URL);
  if (!apiUrl) return next(req); // config not loaded yet (e.g. the env.json request itself)

  const apiOrigin = new URL(apiUrl).origin;
  const target = new URL(req.url, window.location.origin);

  if (target.origin != apiOrigin) return next(req);

  const keycloak = inject(Keycloak);
  if (!keycloak.authenticated) return next(req);

  return from(keycloak.updateToken(30)).pipe(
    switchMap(() => next(req.clone({ setHeaders: { Authorization: `Bearer ${keycloak.token}` } }))),
  );
};
