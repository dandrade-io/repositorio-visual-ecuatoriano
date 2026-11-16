import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { isEmbeddingModel } from '../../assets/services/assets.service';

export const validSearchGuard: CanActivateFn = (route) => {
  const router = inject(Router);

  return isEmbeddingModel(route.queryParamMap.get('model')) ? true : router.createUrlTree(['/']);
};
