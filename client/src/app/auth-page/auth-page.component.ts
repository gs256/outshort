import { Component, effect, inject, Signal, signal } from '@angular/core';
import { ActivatedRoute, Params, Router, RouterLink } from '@angular/router';
import { ButtonModule } from 'primeng/button';
import { CardModule } from 'primeng/card';
import { InputTextModule } from 'primeng/inputtext';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import {
  AbstractControl,
  FormBuilder,
  ReactiveFormsModule,
  ValidationErrors,
  ValidatorFn,
  Validators,
} from '@angular/forms';
import { AuthService } from '../services/api/auth.service';
import { catchError } from 'rxjs';

type FormType = 'sign-in' | 'sign-up';

const DEFAULT_TYPE: FormType = 'sign-up';

@Component({
  selector: 'app-auth-page',
  imports: [
    CardModule,
    ButtonModule,
    InputTextModule,
    RouterLink,
    ReactiveFormsModule,
  ],
  templateUrl: './auth-page.component.html',
  styleUrl: './auth-page.component.scss',
})
export class AuthPageComponent {
  private readonly _route = inject(ActivatedRoute);
  private readonly _router = inject(Router);
  private readonly _fb = inject(FormBuilder);
  private readonly _authService = inject(AuthService);

  public readonly formType = signal<FormType>(DEFAULT_TYPE);
  public readonly formValid = signal(false);

  public readonly form = this._fb.nonNullable.group(
    {
      username: ['', [Validators.required, Validators.minLength(2)]],
      password: ['', [Validators.required, Validators.minLength(6)]],
      repeat: [''],
    },
    { validators: [this.matchPasswordValidator(this.formType)] },
  );

  constructor() {
    this._route.queryParams.pipe(takeUntilDestroyed()).subscribe((params) => {
      this.validateTypeParams(params);
      if (this.getTypeParam(params) === 'sign-in') {
        this.formType.set('sign-in');
      } else {
        this.formType.set('sign-up');
      }
    });
    this.form.valueChanges.pipe(takeUntilDestroyed()).subscribe(() => {
      this.formValid.set(this.form.valid);
    });
    effect(() => {
      this.formType();
      this.form.updateValueAndValidity();
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

  public onSignIn() {
    const formValue = this.form.getRawValue();
    this._authService.signIn(formValue.username, formValue.password).subscribe({
      next: () => {
        this._router.navigate(['/']);
      },
      error: (error: Error) => {
        alert(error.message);
      },
    });
  }

  private matchPasswordValidator(formType: Signal<FormType>): ValidatorFn {
    return (control: AbstractControl): ValidationErrors | null => {
      const password = control.get('password')?.value;
      const repeat = control.get('repeat')?.value;
      if (password === repeat || formType() === 'sign-in') {
        return null;
      }
      return { passwordMismatch: true };
    };
  }
}
