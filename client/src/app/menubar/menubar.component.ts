import { Component, inject } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import { MenuItem } from 'primeng/api';
import { ButtonModule } from 'primeng/button';
import { MenubarModule } from 'primeng/menubar';
import { absoluteRoute } from '../utils';
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
      routerLink: '/app',
    },
    {
      label: 'Dashboard',
      icon: 'pi pi-list',
      routerLink: '/app/dashboard',
    },
    {
      label: 'About',
      icon: 'pi pi-info-circle',
      routerLink: '/app/about',
    },
  ];

  public navigateAuth() {
    this._router.navigate(['app/auth']);
  }

  public navigateUser() {
    this._router.navigate(['app/user']);
  }
}
