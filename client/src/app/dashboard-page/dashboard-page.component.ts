import { Component, inject } from '@angular/core';
import { PageWrapperComponent } from '../page-wrapper/page-wrapper.component';
import { MenubarComponent } from '../menubar/menubar.component';
import { CardModule } from 'primeng/card';
import { ButtonModule } from 'primeng/button';
import { UserStore } from '../store/user.store';
import { Router } from '@angular/router';

@Component({
  selector: 'app-dashboard-page',
  imports: [PageWrapperComponent, MenubarComponent, CardModule, ButtonModule],
  templateUrl: './dashboard-page.component.html',
  styleUrl: './dashboard-page.component.scss',
})
export class DashboardPageComponent {
  private readonly _userStore = inject(UserStore);
  private readonly _router = inject(Router);

  public readonly user = this._userStore.user;

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
}
