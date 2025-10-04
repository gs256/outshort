import { Component, inject, signal } from '@angular/core';
import { ButtonModule } from 'primeng/button';
import { CardModule } from 'primeng/card';
import { TableModule } from 'primeng/table';
import { LinksStore } from '../../store/links.store';
import { Link } from '../../models/link';
import { getShortUrl } from '../../utils';
import { LinkModalComponent } from '../../link-modal/link-modal.component';
import { Clipboard } from '@angular/cdk/clipboard';

@Component({
  selector: 'app-dashboard-links',
  imports: [CardModule, ButtonModule, TableModule, LinkModalComponent],
  templateUrl: './dashboard-links.component.html',
  styleUrl: './dashboard-links.component.css',
})
export class DashboardLinksComponent {
  private readonly _linksStore = inject(LinksStore);
  private readonly _clipboard = inject(Clipboard);

  public readonly editorVisible = signal(false);
  public readonly links = this._linksStore.links;
  public readonly getShortUrl = getShortUrl;

  public onRowClicked(link: Link) {
    this._linksStore.setDraft(link.id);
    this.editorVisible.set(true);
  }

  public onCopyLink(event: Event, link: Link) {
    event.stopPropagation();
    this._clipboard.copy(getShortUrl(link.alias));
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
