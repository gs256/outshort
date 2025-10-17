import { Injectable } from '@angular/core';

export type AppTheme = 'light' | 'dark';

const THEME_KEY = 'theme';

@Injectable({
  providedIn: 'root',
})
export class AppThemeService {
  public setSavedTheme() {
    this.setTheme(this.getCurrentTheme());
  }

  public getCurrentTheme(): AppTheme {
    const themeValue = localStorage.getItem(THEME_KEY);
    const theme = themeValue === 'dark' ? 'dark' : 'light';
    localStorage.setItem(THEME_KEY, theme);
    return theme;
  }

  public setTheme(theme: AppTheme) {
    const html = document.querySelector('html');
    if (theme === 'dark') {
      html?.classList.add('dark-theme');
    } else {
      html?.classList.remove('dark-theme');
    }
    localStorage.setItem(THEME_KEY, theme);
  }

  public toggleTheme() {
    if (this.getCurrentTheme() === 'dark') {
      this.setTheme('light');
    } else {
      this.setTheme('dark');
    }
  }
}
