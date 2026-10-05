import { Routes } from '@angular/router';
import { HomeLayout } from './home/layouts/home-layout/home-layout';
import { ContributeLayout } from './contribute/layouts/contribute-layout/contribute-layout';
import { RegisterAlbumCoverPage } from './album-cover/pages/register-album-cover-page/register-album-cover-page';

export const routes: Routes = [
  {
    path: 'contribute',
    component: ContributeLayout,
    children: [
      { path: '', redirectTo: 'album-cover', pathMatch: 'full' },
      { path: 'album-cover', component: RegisterAlbumCoverPage },
    ],
  },
  {
    path: '',
    pathMatch: 'full',
    component: HomeLayout,
  },
];
