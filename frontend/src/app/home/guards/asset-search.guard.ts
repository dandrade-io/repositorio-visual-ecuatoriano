import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { ArtworkSearchService } from '../services/artwork-search-service';

export const assetSearchGuard: CanActivateFn = (route) => {
  const router = inject(Router);
  const searchService = inject(ArtworkSearchService);

  const isImageMode = route.queryParamMap.get('mode') === 'image';

  if (isImageMode && !searchService.asset()) {
    return router.createUrlTree(['/']);
  }

  return true;
};
