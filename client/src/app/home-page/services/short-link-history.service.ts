import { Injectable, signal } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class ShortLinkHistoryService {
  private readonly _history = signal<
    Array<{ original: string; alias: string }>
  >([]);

  // TODO: make readonly
  public readonly records = this._history;

  public add(originalUrl: string, alias: string) {
    this._history.set([
      { original: originalUrl, alias: alias },
      ...this._history(),
    ]);
  }
}
