import { Routes } from '@angular/router';
import { HomePageComponent } from './home-page/home-page.component';
import { RedirectComponent } from './redirect/redirect.component';
import { ROUTES } from './constants';
import { AuthPageComponent } from './auth-page/auth-page.component';
import { DashboardPageComponent } from './dashboard-page/dashboard-page.component';
import { AboutPageComponent } from './about-page/about-page.component';
import { UserPageComponent } from './user-page/user-page.component';

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
    component: DashboardPageComponent,
  },
  {
    path: ROUTES.about,
    component: AboutPageComponent,
  },
  {
    path: ROUTES.auth,
    component: AuthPageComponent,
  },
  {
    path: ROUTES.user,
    component: UserPageComponent,
  },
  {
    path: ':alias',
    component: RedirectComponent,
  },
];
