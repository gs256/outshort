import {
  patchState,
  signalStore,
  withHooks,
  withMethods,
  withState,
} from '@ngrx/signals';
import { User } from '../models/user';
import { inject } from '@angular/core';
import { AuthService } from '../services/api/auth.service';
import { pipe, switchMap, tap } from 'rxjs';
import { rxMethod } from '@ngrx/signals/rxjs-interop';
import { tapResponse } from '@ngrx/operators';

type UserState = {
  user: User | undefined;
  isLoading: boolean;
};

const initialState: UserState = {
  user: undefined,
  isLoading: false,
};

export const UserStore = signalStore(
  { providedIn: 'root' },
  withState(initialState),
  withMethods((store, authService = inject(AuthService)) => ({
    load: rxMethod(
      pipe(
        tap(() => patchState(store, { isLoading: true })),
        switchMap((query) => {
          return authService.getUserInfo().pipe(
            tapResponse({
              next: (user) => patchState(store, { user }),
              error: () => {},
              finalize: () => patchState(store, { isLoading: false }),
            }),
          );
        }),
      ),
    ),
  })),
  withHooks({
    onInit(store) {
      store.load({});
    },
  }),
);
