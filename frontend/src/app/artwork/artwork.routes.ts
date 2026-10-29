import { Routes } from '@angular/router';
import { ArtworkLayout } from './layouts/artwork-layout/artwork-layout';
import { AlbumCoverPage } from '../album-cover/pages/album-cover-page/album-cover-page';

export const artworkRoutes: Routes = [
  {
    path: '',
    component: ArtworkLayout,
    children: [
      {
        path: 'album-cover/:id',
        component: AlbumCoverPage,
      },
    ],
  },
];
