import { Component, inject, signal } from '@angular/core';
import { ActivatedRoute, Params, Router, RouterLink } from '@angular/router';
import { ButtonModule } from 'primeng/button';
import { CardModule } from 'primeng/card';
import { InputTextModule } from 'primeng/inputtext';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

type FormType = 'sign-in' | 'sign-up';

const DEFAULT_TYPE: FormType = 'sign-up';

@Component({
  selector: 'app-auth-page',
  imports: [CardModule, ButtonModule, InputTextModule, RouterLink],
  templateUrl: './auth-page.component.html',
  styleUrl: './auth-page.component.scss',
})
export class AuthPageComponent {
  private readonly _route = inject(ActivatedRoute);
  private readonly _router = inject(Router);

  public readonly formType = signal<FormType>(DEFAULT_TYPE);

  constructor() {
    this._route.queryParams.pipe(takeUntilDestroyed()).subscribe((params) => {
      this.validateTypeParams(params);
      if (this.getTypeParam(params) === 'sign-in') {
        this.formType.set('sign-in');
      } else {
        this.formType.set('sign-up');
      }
    });
  }

  private validateTypeParams(params: Params) {
    const typeParam = this.getTypeParam(params);
    if (typeParam !== 'sign-in' && typeParam !== 'sign-up') {
      this._router.navigate([], {
        relativeTo: this._route,
        replaceUrl: true,
        queryParams: { type: DEFAULT_TYPE },
        queryParamsHandling: 'merge',
      });
    }
  }

  private getTypeParam(params: Params) {
    return params['type'];
  }
}
