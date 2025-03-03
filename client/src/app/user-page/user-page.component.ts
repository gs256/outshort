import { Component, effect, inject } from '@angular/core';
import { PageWrapperComponent } from '../page-wrapper/page-wrapper.component';
import { MenubarComponent } from '../menubar/menubar.component';
import { CardModule } from 'primeng/card';
import { UserStore } from '../store/user.store';
import { Router } from '@angular/router';

@Component({
  selector: 'app-user-page',
  imports: [PageWrapperComponent, MenubarComponent, CardModule],
  templateUrl: './user-page.component.html',
  styleUrl: './user-page.component.scss',
})
export class UserPageComponent {
  private readonly _userStore = inject(UserStore);
  private readonly _router = inject(Router);

  public readonly user = this._userStore.user;

  constructor() {
    effect(() => {
      if (this._userStore.user() === undefined) {
        this._router.navigate(['app']);
      }
    });
  }

  public signOut() {
    this._userStore.signOut({});
  }
}
