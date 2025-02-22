import { Component, inject, signal } from '@angular/core';
import { ApiService } from '../services/api/api.service';
import { finalize } from 'rxjs';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-home-page',
  imports: [FormsModule],
  templateUrl: './home-page.component.html',
  styleUrl: './home-page.component.scss',
})
export class HomePageComponent {
  private readonly api = inject(ApiService);

  public readonly processing = signal(false);

  public readonly shortLink = signal('');

  public readonly originalUrl = signal('');

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
