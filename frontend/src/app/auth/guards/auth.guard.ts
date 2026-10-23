import { inject } from '@angular/core';
import {
  ActivatedRouteSnapshot,
  CanActivateFn,
  Router,
  RouterStateSnapshot,
  UrlTree,
} from '@angular/router';
import { AuthGuardData, createAuthGuard } from 'keycloak-angular';

const isAccessAllowed = async (
  route: ActivatedRouteSnapshot,
  _state: RouterStateSnapshot,
  authData: AuthGuardData,
): Promise<boolean | UrlTree> => {
  const { authenticated, grantedRoles } = authData;
  const router = inject(Router);

  if (!authenticated) return router.parseUrl('/');

  const requiredRole = route.data['role'];
  if (!requiredRole) return true;

  const hasRequiredRole = Object.values(grantedRoles.resourceRoles).some((roles) =>
    roles.includes(requiredRole),
  );

  return hasRequiredRole ? true : router.parseUrl('/');
};

export const authGuard = createAuthGuard<CanActivateFn>(isAccessAllowed);
