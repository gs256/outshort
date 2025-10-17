import { Component, inject } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { AppThemeService } from './services/app-theme.service';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css',
})
export class AppComponent {
  private readonly _appTheme = inject(AppThemeService);

  constructor() {
    this._appTheme.setSavedTheme();
  }
}
