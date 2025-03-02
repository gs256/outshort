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
import { storage } from '../storage';

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
  withMethods((store, authService = inject(AuthService)) => {
    const load = rxMethod(
      pipe(
        tap(() => patchState(store, { isLoading: true })),
        switchMap(() => {
          return authService.getUserInfo().pipe(
            tapResponse({
              next: (user) => patchState(store, { user }),
              error: () => {},
              finalize: () => patchState(store, { isLoading: false }),
            }),
          );
        }),
      ),
    );

    const signIn = rxMethod<{ username: string; password: string }>(
      pipe(
        tap(() => patchState(store, { isLoading: true })),
        switchMap((params) => {
          return authService.signIn(params.username, params.password).pipe(
            tapResponse({
              next: (authToken) => {
                storage.authToken = authToken;
                load({});
              },
              error: () => {},
              finalize: () => patchState(store, { isLoading: false }),
            }),
          );
        }),
      ),
    );

    const signOut = rxMethod(
      pipe(
        tap(() => patchState(store, { isLoading: true })),
        switchMap(() => {
          return authService.signOut().pipe(
            tapResponse({
              next: () => {
                patchState(store, { user: undefined });
              },
              error: () => {},
              finalize: () => patchState(store, { isLoading: false }),
            }),
          );
        }),
      ),
    );

    return { load, signIn, signOut };
  }),
  withHooks({
    onInit(store) {
      store.load({});
    },
  }),
);
