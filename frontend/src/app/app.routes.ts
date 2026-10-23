import { Routes } from '@angular/router';
import { HomeLayout } from './home/layouts/home-layout/home-layout';
import { ContributeLayout } from './contribute/layouts/contribute-layout/contribute-layout';
import { RegisterAlbumCoverPage } from './album-cover/pages/register-album-cover-page/register-album-cover-page';
import { authGuard } from './auth/guards/auth.guard';

export const routes: Routes = [
  {
    path: '',
    component: HomeLayout,
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
