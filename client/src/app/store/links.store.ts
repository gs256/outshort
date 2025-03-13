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
import { LinkUpsert } from '../models/link-upsert';

export const LinksStore = signalStore(
  withState({
    links: [] as Link[],
    draftId: null as string | null,
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

    const createLink = rxMethod<LinkUpsert>(
      pipe(
        tap(() => patchState(store, { isLoading: true })),
        switchMap((body) => {
          return linksService.createLink(body).pipe(
            tapResponse({
              next: () => load({}),
              error: () => {},
              finalize: () => patchState(store, { isLoading: false }),
            }),
          );
        }),
      ),
    );

    const updateLink = rxMethod<{ id: string; body: LinkUpsert }>(
      pipe(
        tap(() => patchState(store, { isLoading: true })),
        switchMap(({ id, body }) => {
          return linksService.updateLink(id, body).pipe(
            tapResponse({
              next: () => load({}),
              error: () => {},
              finalize: () => patchState(store, { isLoading: false }),
            }),
          );
        }),
      ),
    );

    const setDraft = (draftId?: string | null) => {
      patchState(store, { draftId: draftId ?? null });
    };

    const findLink = (linkId?: string | null) => {
      return linkId
        ? (store.links().find((item) => item.id === linkId) ?? null)
        : null;
    };

    return { load, createLink, setDraft, findLink, updateLink };
  }),
  withHooks({
    onInit(store) {
      store.load({});
    },
  }),
);
