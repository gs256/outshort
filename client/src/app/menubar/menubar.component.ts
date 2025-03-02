import { Component, inject } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import { MenuItem } from 'primeng/api';
import { ButtonModule } from 'primeng/button';
import { MenubarModule } from 'primeng/menubar';
import { absoluteRoute } from '../utils';
import { ROUTES } from '../constants';
import { UserStore } from '../store/user.store';
import { ChipModule } from 'primeng/chip';

@Component({
  selector: 'app-menubar',
  imports: [MenubarModule, RouterLink, ButtonModule, ChipModule],
  templateUrl: './menubar.component.html',
  styleUrl: './menubar.component.scss',
})
export class MenubarComponent {
  private readonly _router = inject(Router);
  private readonly _userStore = inject(UserStore);

  public readonly loading = this._userStore.isLoading;
  public readonly user = this._userStore.user;

  public readonly menuItems: MenuItem[] = [
    {
      label: 'Home',
      icon: 'pi pi-home',
      routerLink: absoluteRoute(ROUTES.home),
    },
    {
      label: 'Dashboard',
      icon: 'pi pi-list',
      routerLink: absoluteRoute(ROUTES.dashboard),
    },
    {
      label: 'About',
      icon: 'pi pi-info-circle',
      routerLink: absoluteRoute(ROUTES.about),
    },
  ];

  public navigateAuth() {
    this._router.navigate([ROUTES.auth]);
  }
}
