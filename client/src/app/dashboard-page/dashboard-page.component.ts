import { Component, inject, signal } from '@angular/core';
import { PageWrapperComponent } from '../page-wrapper/page-wrapper.component';
import { MenubarComponent } from '../menubar/menubar.component';
import { CardModule } from 'primeng/card';
import { ButtonModule } from 'primeng/button';
import { UserStore } from '../store/user.store';
import { Router } from '@angular/router';
import { Link } from '../models/link';
import { TableModule } from 'primeng/table';
import { LinksStore } from '../store/links.store';
import { getShortUrl } from '../utils';

@Component({
  selector: 'app-dashboard-page',
  imports: [
    PageWrapperComponent,
    MenubarComponent,
    CardModule,
    ButtonModule,
    TableModule,
  ],
  templateUrl: './dashboard-page.component.html',
  styleUrl: './dashboard-page.component.scss',
  providers: [LinksStore],
})
export class DashboardPageComponent {
  private readonly _userStore = inject(UserStore);
  private readonly _router = inject(Router);
  private readonly _linksStore = inject(LinksStore);

  public readonly user = this._userStore.user;
  public readonly links = this._linksStore.links;
  public readonly getShortUrl = getShortUrl;

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

  public getExpirationString(link: Link) {
    const creation = new Date(link.creationDate).getTime();
    const lifetimeMs = link.lifetime * 1000;
    if (lifetimeMs === 0) {
      return 'Never';
    }
    const expirationDate = new Date(creation + lifetimeMs);
    const now = Date.now();
    const timeLeftSec = Math.floor((expirationDate.getTime() - now) / 1000);
    if (timeLeftSec < 0) {
      return 'Expired';
    }
    const timeString = this.convertSeconds(timeLeftSec);
    return `In ${timeString}`;
  }

  public lifetimeToString(lifetimeSec: number) {
    switch (lifetimeSec) {
      case 0:
        return 'Forever';
      default:
        return this.convertSeconds(lifetimeSec);
    }
  }

  public onRowClicked(link: Link) {
    console.log('row clicked', link.alias);
  }

  private convertSeconds(seconds: number): string {
    const units = [
      { label: 'year', seconds: 31536000 },
      { label: 'month', seconds: 2592000 },
      { label: 'week', seconds: 604800 },
      { label: 'day', seconds: 86400 },
      { label: 'hour', seconds: 3600 },
      { label: 'minute', seconds: 60 },
    ];
    for (const unit of units) {
      const value = Math.floor(seconds / unit.seconds);
      if (value >= 1) {
        return `${value} ${unit.label}${value > 1 ? 's' : ''}`;
      }
    }
    return `${seconds} seconds`;
  }
}
