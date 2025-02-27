import { Component, inject, signal } from '@angular/core';
import { ApiService } from '../services/api/api.service';
import { finalize } from 'rxjs';
import { FormsModule } from '@angular/forms';
import { ButtonModule } from 'primeng/button';
import { MenubarModule } from 'primeng/menubar';
import { MenuItem, MessageService } from 'primeng/api';
import { CardModule } from 'primeng/card';
import { InputTextModule } from 'primeng/inputtext';
import { TabsModule } from 'primeng/tabs';
import { ClipboardModule } from '@angular/cdk/clipboard';
import { Clipboard } from '@angular/cdk/clipboard';
import { ToastModule } from 'primeng/toast';

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
  ],
  templateUrl: './home-page.component.html',
  styleUrl: './home-page.component.scss',
  providers: [MessageService],
})
export class HomePageComponent {
  private readonly _api = inject(ApiService);
  private readonly _clipboard = inject(Clipboard);
  private readonly _messageService = inject(MessageService);

  public readonly processing = signal(false);
  public readonly shortLink = signal('');
  public readonly originalUrl = signal('');
  public readonly menuItems: MenuItem[] = [
    {
      label: 'Home',
      icon: 'pi pi-home',
      command: () => {
        console.log('Home');
      },
    },
    {
      label: 'Dashboard',
      icon: 'pi pi-list',
      command: () => {
        console.log('Dashboard');
      },
    },
    {
      label: 'About',
      icon: 'pi pi-info-circle',
      command: () => {
        console.log('About');
      },
    },
  ];

  public onShortenClicked() {
    if (this.processing()) {
      return;
    }
    if (this.originalUrl().trim().length == 0) {
      alert('Enter your url');
      return;
    }
    this.processing.set(true);
    this._api
      .shorten(this.originalUrl())
      .pipe(
        finalize(() => {
          this.processing.set(false);
        })
      )
      .subscribe({
        next: (result) => {
          this.shortLink.set(this.getShortUrl(result));
        },
        error: (error: Error) => {
          alert(error.message);
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

  private getShortUrl(alias: string) {
    return `${window.location.origin}/${alias}`;
  }
}
