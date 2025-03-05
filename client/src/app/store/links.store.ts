import {
  patchState,
  signalStore,
  withHooks,
  withMethods,
  withState,
} from '@ngrx/signals';
import { inject } from '@angular/core';
import { pipe, switchMap, tap } from 'rxjs';
import { rxMethod } from '@ngrx/signals/rxjs-interop';
import { tapResponse } from '@ngrx/operators';
import { LinksService } from '../services/api/links.service';
import { Link } from '../models/link';

export const LinksStore = signalStore(
  withState({
    links: [] as Link[],
    isLoading: false,
  }),
  withMethods((store, linksService = inject(LinksService)) => {
    const load = rxMethod(
      pipe(
        tap(() => patchState(store, { isLoading: true })),
        switchMap(() => {
          return linksService.getAllLinks().pipe(
            tapResponse({
              next: (links) => patchState(store, { links }),
              error: () => {},
              finalize: () => patchState(store, { isLoading: false }),
            }),
          );
        }),
      ),
    );

    return { load };
  }),
  withHooks({
    onInit(store) {
      store.load({});
    },
  }),
);
