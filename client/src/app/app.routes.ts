import { Routes } from '@angular/router';
import { HomePageComponent } from './home-page/home-page.component';
import { RedirectComponent } from './redirect/redirect.component';
import { ROUTES } from './constants';

export const routes: Routes = [
  {
    path: '',
    redirectTo: ROUTES.home,
    pathMatch: 'full',
  },
  {
    path: ROUTES.home,
    component: HomePageComponent,
  },
  {
    path: ROUTES.dashboard,
    component: HomePageComponent,
  },
  {
    path: ROUTES.about,
    component: HomePageComponent,
  },
  {
    path: ROUTES.auth,
    component: HomePageComponent,
  },
  {
    path: ':alias',
    component: RedirectComponent,
  },
];
