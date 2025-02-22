import { Component, inject } from '@angular/core';
import { ApiService } from '../services/api/api.service';

@Component({
  selector: 'app-home-page',
  imports: [],
  templateUrl: './home-page.component.html',
  styleUrl: './home-page.component.scss',
})
export class HomePageComponent {
  private readonly api = inject(ApiService);

  public onShortenClicked() {
    this.api.shorten('https://example.com').subscribe((result) => {
      console.log('result:', result);
    });
  }
}
