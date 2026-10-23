import { Routes } from '@angular/router';
import { ContributeLayout } from './layouts/contribute-layout/contribute-layout';
import { RegisterAlbumCoverPage } from '../album-cover/pages/register-album-cover-page/register-album-cover-page';

export const contributeRoutes: Routes = [
  {
    path: '',
    component: ContributeLayout,
    title: 'Registro de Portada',
    children: [
      { path: 'album-cover', component: RegisterAlbumCoverPage },
      { path: '**', redirectTo: 'album-cover' },
    ],
  },
];

export default contributeRoutes;
