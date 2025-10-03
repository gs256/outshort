import { Component, inject, signal } from '@angular/core';
import { PageWrapperComponent } from '../../page-wrapper/page-wrapper.component';
import { MenubarComponent } from '../../menubar/menubar.component';
import { CardModule } from 'primeng/card';
import { ButtonModule } from 'primeng/button';
import { UserStore } from '../../store/user.store';
import { ActivatedRoute, Router, RouterOutlet } from '@angular/router';
import { Link } from '../../models/link';
import { TableModule } from 'primeng/table';
import { LinksStore } from '../../store/links.store';
import { MenuModule } from 'primeng/menu';
import { MenuItem } from 'primeng/api';
import { MenubarModule } from 'primeng/menubar';
import { getResolvedUser } from '../../utils/get-resolved-user';

@Component({
  selector: 'app-dashboard-page',
  imports: [
    PageWrapperComponent,
    MenubarComponent,
    CardModule,
    ButtonModule,
    TableModule,
    MenuModule,
    MenubarModule,
    RouterOutlet,
  ],
  templateUrl: './dashboard-page.component.html',
  styleUrl: './dashboard-page.component.scss',
  providers: [LinksStore],
})
export class DashboardPageComponent {
  private readonly _route = inject(ActivatedRoute);
  private readonly _router = inject(Router);

  public readonly user = getResolvedUser(this._route);
  public readonly menuItems: MenuItem[] = [
    {
      label: 'Links',
      icon: 'pi pi-link',
      routerLink: 'links',
    },
    {
      label: 'Pages',
      icon: 'pi pi-book',
      routerLink: 'pages',
    },
    {
      label: 'Analytics',
      icon: 'pi pi-chart-bar',
      routerLink: 'analytics',
    },
    {
      label: 'Settings',
      icon: 'pi pi-cog',
      routerLink: 'settings',
    },
  ];

  public navigateSignIn() {
    this._router.navigate(['app/auth'], {
      queryParams: { type: 'sign-in' },
    });
  }

  public navigateSignUp() {
    this._router.navigate(['app/auth'], {
      queryParams: { type: 'sign-up' },
    });
  }

  public getExpirationDate(link: Link): string {
    const creation = new Date(link.creationDate).getTime();
    const lifetimeMs = link.lifetime * 1000;
    if (lifetimeMs === 0) {
      return 'Never';
    }
    const expirationDate = new Date(creation + lifetimeMs);
    const year = expirationDate.getFullYear();
    const month = String(expirationDate.getMonth() + 1).padStart(2, '0');
    const day = String(expirationDate.getDate()).padStart(2, '0');
    return `${year}.${month}.${day}`;
  }
}
