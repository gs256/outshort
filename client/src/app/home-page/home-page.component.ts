import { Component, inject, signal } from '@angular/core';
import { ApiService } from '../services/api/api.service';
import { finalize } from 'rxjs';
import { FormsModule } from '@angular/forms';
import { ButtonModule } from 'primeng/button';
import { MenubarModule } from 'primeng/menubar';
import { MenuItem } from 'primeng/api';
import { CardModule } from 'primeng/card';

@Component({
  selector: 'app-home-page',
  imports: [FormsModule, ButtonModule, MenubarModule, CardModule],
  templateUrl: './home-page.component.html',
  styleUrl: './home-page.component.scss',
})
export class HomePageComponent {
  private readonly api = inject(ApiService);
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
    this.api
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

  private getShortUrl(alias: string) {
    return `${window.location.origin}/${alias}`;
  }
}
