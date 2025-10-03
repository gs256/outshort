import {
  ActivatedRouteSnapshot,
  ResolveFn,
  RouterStateSnapshot,
} from '@angular/router';
import { User } from '../models/user';
import { UserStore } from '../store/user.store';
import { inject } from '@angular/core';
import { toObservable } from '@angular/core/rxjs-interop';
import { combineLatest, filter, map, skip } from 'rxjs';

export const userResolver: ResolveFn<User | undefined> = (
  route: ActivatedRouteSnapshot,
  state: RouterStateSnapshot,
) => {
  const userStore = inject(UserStore);

  const isLoading$ = toObservable(userStore.isLoading);
  const user$ = toObservable(userStore.user);

  // User already resolved
  if (userStore.initialized() && !userStore.isLoading()) {
    return userStore.user();
  }

  return combineLatest([isLoading$, user$]).pipe(
    filter(([loading, _]) => !loading),
    skip(1),
    map(([_, user]) => user),
  );
};
