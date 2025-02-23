import { Routes } from '@angular/router';
import { HomePageComponent } from './home-page/home-page.component';
import { RedirectComponent } from './redirect/redirect.component';

export const routes: Routes = [
  {
    path: ':alias',
    component: RedirectComponent,
  },
  {
    path: '',
    component: HomePageComponent,
  },
];
