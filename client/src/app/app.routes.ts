import { Routes } from '@angular/router';
import { HomePageComponent } from './home-page/home-page.component';
import { RedirectComponent } from './redirect/redirect.component';
import { AuthPageComponent } from './auth-page/auth-page.component';
import { DashboardPageComponent } from './dashboard/dashboard-page/dashboard-page.component';
import { AboutPageComponent } from './about-page/about-page.component';
import { UserPageComponent } from './user-page/user-page.component';
import { DashboardLinksComponent } from './dashboard/dashboard-links/dashboard-links.component';
import { DashboardPagesComponent } from './dashboard/dashboard-pages/dashboard-pages.component';
import { DashboardAnalyticsComponent } from './dashboard/dashboard-analytics/dashboard-analytics.component';
import { DashboardSettingsComponent } from './dashboard/dashboard-settings/dashboard-settings.component';

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
    children: [
      {
        path: 'links',
        component: DashboardLinksComponent,
      },
      {
        path: 'pages',
        component: DashboardPagesComponent,
      },
      {
        path: 'analytics',
        component: DashboardAnalyticsComponent,
      },
      {
        path: 'settings',
        component: DashboardSettingsComponent,
      },
      {
        path: '**',
        redirectTo: 'links',
      },
    ],
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
