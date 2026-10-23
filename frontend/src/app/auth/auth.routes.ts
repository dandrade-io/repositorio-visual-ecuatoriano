import { Routes } from '@angular/router';
import { AuthLayout } from './layouts/auth.layout/auth-layout';
import { CreateAccountPage } from './pages/register/register-page';

export const authRoutes: Routes = [
  {
    path: '',
    component: AuthLayout,
    children: [
      { path: 'register', component: CreateAccountPage },
      { path: '**', redirectTo: window.location.origin },
    ],
  },
];

export default authRoutes;
