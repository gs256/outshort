import { Component, computed, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ButtonModule } from 'primeng/button';
import { MenubarModule } from 'primeng/menubar';
import { MessageService } from 'primeng/api';
import { CardModule } from 'primeng/card';
import { InputTextModule } from 'primeng/inputtext';
import { TabsModule } from 'primeng/tabs';
import { ClipboardModule } from '@angular/cdk/clipboard';
import { Clipboard } from '@angular/cdk/clipboard';
import { ToastModule } from 'primeng/toast';
import { Router, RouterLink } from '@angular/router';
import { TableModule } from 'primeng/table';
import { MenubarComponent } from '../menubar/menubar.component';
import { PageWrapperComponent } from '../page-wrapper/page-wrapper.component';
import { UserStore } from '../store/user.store';
import { getShortUrl } from '../utils';
import { SelectButtonModule } from 'primeng/selectbutton';
import { PageContentWrapperComponent } from '../page-content-wrapper/page-content-wrapper.component';
import { QuickLinksStore } from '../store/quick-links.store';

@Component({
  selector: 'app-home-page',
  imports: [
    FormsModule,
    ButtonModule,
    MenubarModule,
    CardModule,
    InputTextModule,
    TabsModule,
    ClipboardModule,
    ToastModule,
    TableModule,
    MenubarComponent,
    PageWrapperComponent,
    RouterLink,
    SelectButtonModule,
    PageContentWrapperComponent,
  ],
  templateUrl: './home-page.component.html',
  styleUrl: './home-page.component.css',
  providers: [MessageService],
})
export class HomePageComponent {
  private readonly _clipboard = inject(Clipboard);
  private readonly _messageService = inject(MessageService);
  private readonly _router = inject(Router);
  private readonly _userStore = inject(UserStore);
  private readonly _message = inject(MessageService);
  private readonly _store = inject(QuickLinksStore);

  public readonly isLoadgin = this._store.isLoading;
  public readonly shortLink = this._store.link;
  public readonly originalUrl = signal('');
  public readonly quickLinkLifetime = signal('1h');
  public readonly user = this._userStore.user;
  public readonly getShortUrl = getShortUrl;

  public readonly lifetimeOptions = [
    { label: '1 hour', value: '1h' },
    { label: '1 day', value: '1d' },
    { label: '1 week', value: '1w' },
  ];

  public readonly shortened = computed(() => this.shortLink() !== null);

  public async onShortenClicked() {
    if (this.isLoadgin()) {
      return;
    }
    const originalUrl = this.originalUrl().trim();
    if (originalUrl.length == 0) {
      return;
    }
    const alias = await this._store.shorten(this.originalUrl());
    if (alias === null) {
      this._message.add({
        severity: 'error',
        summary: 'Error',
        detail: 'Invalid url',
      });
    }
  }

  public copyShortLink() {
    const shortLink = this.shortLink();
    if (!shortLink) {
      return;
    }
    this._clipboard.copy(shortLink);
    this._messageService.add({
      summary: 'Copied to clipboard',
      detail: shortLink,
      severity: 'success',
    });
  }

  public navigateAuth() {
    this._router.navigate(['app/auth']);
  }

  public navigateDashboard() {
    this._router.navigate(['app/dashboard']);
  }
}
