import { Routes } from '@angular/router';
import { HomeLayout } from './home/layouts/home-layout/home-layout';
import { authGuard } from './auth/guards/auth.guard';
import { artworkRoutes } from './artwork/artwork.routes';

export const routes: Routes = [
  {
    path: '',
    component: HomeLayout,
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
