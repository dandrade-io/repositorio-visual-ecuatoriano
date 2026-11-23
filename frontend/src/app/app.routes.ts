import { Routes } from '@angular/router';
import { HomeLayout } from './home/layouts/home-layout/home-layout';
import { authGuard } from './auth/guards/auth.guard';
import { artworkRoutes } from './artwork/artwork.routes';
import { validSearchGuard } from './home/guards/valid-search.guard';
import { HomeResultsPage } from './home/pages/home-results-page/home-results-page';
import { assetSearchGuard } from './home/guards/asset-search.guard';

export const routes: Routes = [
  {
    path: '',
    component: HomeLayout,
    children: [
      { path: '', component: HomeResultsPage, pathMatch: 'full' },
      { path: 'search', component: HomeResultsPage, canActivate: [validSearchGuard, assetSearchGuard] },
    ],
  },
  {
    path: 'artwork',
    children: artworkRoutes,
  },
  {
    path: 'auth',
    loadChildren: () => import('./auth/auth.routes'),
  },
  {
    path: 'contribute',
    canActivate: [authGuard],
    loadChildren: () => import('./contribute/contribute.routes'),
  },
];
