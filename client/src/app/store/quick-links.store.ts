import { inject } from '@angular/core';
import { patchState, signalStore, withMethods, withState } from '@ngrx/signals';
import { ApiService } from '../services/api/api.service';
import { catchError, finalize, lastValueFrom, map, of, tap } from 'rxjs';
import { getShortUrl } from '../utils';

interface QuickLinksState {
  link: string | null;
  isLoading: boolean;
}

export const QuickLinksStore = signalStore(
  { providedIn: 'root' },
  withState<QuickLinksState>({
    link: null,
    isLoading: false,
  }),
  withMethods((store, api = inject(ApiService)) => {
    const shorten = async (originalUrl: string) => {
      patchState(store, { isLoading: true });
      return await lastValueFrom(
        api.quickShorten(originalUrl).pipe(
          map((value) => value),
          tap((alias) => {
            patchState(store, { link: getShortUrl(alias) });
          }),
          catchError(() => of(null)),
          finalize(() => {
            patchState(store, { isLoading: false });
          }),
        ),
      );
    };

    return { shorten };
  }),
);
