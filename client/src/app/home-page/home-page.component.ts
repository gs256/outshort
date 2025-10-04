import { Component, computed, inject, signal } from '@angular/core';
import { ApiService } from '../services/api/api.service';
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
import { ShortLinkHistoryService } from './services/short-link-history.service';
import { TableModule } from 'primeng/table';
import { MenubarComponent } from '../menubar/menubar.component';
import { PageWrapperComponent } from '../page-wrapper/page-wrapper.component';
import { UserStore } from '../store/user.store';
import { getShortUrl } from '../utils';

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
  ],
  templateUrl: './home-page.component.html',
  styleUrl: './home-page.component.css',
  providers: [MessageService],
})
export class HomePageComponent {
  private readonly _api = inject(ApiService);
  private readonly _clipboard = inject(Clipboard);
  private readonly _messageService = inject(MessageService);
  private readonly _historyService = inject(ShortLinkHistoryService);
  private readonly _router = inject(Router);
  private readonly _userStore = inject(UserStore);

  public readonly processing = signal(false);
  public readonly shortLink = signal('');
  public readonly originalUrl = signal('');
  public readonly history = this._historyService.records;
  public readonly user = this._userStore.user;
  public readonly getShortUrl = getShortUrl;

  public readonly shortened = computed(
    () => this.shortLink().trim().length > 0,
  );

  public onShortenClicked() {
    if (this.processing()) {
      return;
    }
    const originalUrl = this.originalUrl().trim();
    if (originalUrl.length == 0) {
      alert('Enter your url');
      return;
    }
    this.processing.set(true);
    this._api.quickShorten(this.originalUrl()).subscribe({
      next: (alias) => {
        this.shortLink.set(getShortUrl(alias));
        this._historyService.add(originalUrl, alias);
        this.processing.set(false);
      },
      error: (error: Error) => {
        alert(error.message);
        this.processing.set(false);
      },
    });
  }

  public copyShortLink() {
    if (this.shortLink().length === 0) {
      return;
    }
    this._clipboard.copy(this.shortLink());
    this._messageService.add({
      summary: 'Copied to clipboard',
      detail: this.shortLink(),
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
