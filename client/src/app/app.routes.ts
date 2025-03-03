import { Routes } from '@angular/router';
import { HomePageComponent } from './home-page/home-page.component';
import { RedirectComponent } from './redirect/redirect.component';
import { AuthPageComponent } from './auth-page/auth-page.component';
import { DashboardPageComponent } from './dashboard-page/dashboard-page.component';
import { AboutPageComponent } from './about-page/about-page.component';
import { UserPageComponent } from './user-page/user-page.component';

export const routes: Routes = [
  {
    path: '',
    redirectTo: 'app',
    pathMatch: 'full',
  },
  {
    path: 'app',
    component: HomePageComponent,
  },
  {
    path: 'app/dashboard',
    component: DashboardPageComponent,
  },
  {
    path: 'app/about',
    component: AboutPageComponent,
  },
  {
    path: 'app/auth',
    component: AuthPageComponent,
  },
  {
    path: 'app/user',
    component: UserPageComponent,
  },
  {
    path: ':alias',
    component: RedirectComponent,
  },
];
