import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { catchError, map, Observable, throwError } from 'rxjs';
import { API_URL } from '../../constants';
import { Link } from '../../models/link';
import { storage } from '../../storage';
import { getErrorResponseMessage } from '../../utils';
import { LinkUpsert } from '../../models/link-upsert';

@Injectable({
  providedIn: 'root',
})
export class LinksService {
  private readonly _http = inject(HttpClient);

  public getAllLinks(): Observable<Link[]> {
    const authToken = storage.authToken;
    if (authToken.length === 0) {
      return throwError(() => new Error('No auth token'));
    }
    return this._http
      .get(`${API_URL}/api/v1/links/all`, {
        headers: { Authorization: `Bearer ${authToken}` },
      })
      .pipe(
        catchError((error: HttpErrorResponse) => {
          const message = getErrorResponseMessage(error);
          if (message) throw new Error(message);
          throw new Error('Unknown error occured');
        }),
        map((res) => {
          const links = res as Link[];
          return links;
        }),
      );
  }

  public createLink(body: LinkUpsert): Observable<Link> {
    const authToken = storage.authToken;
    if (authToken.length === 0) {
      return throwError(() => new Error('No auth token'));
    }
    return this._http
      .post(`${API_URL}/api/v1/links/create`, body, {
        headers: { Authorization: `Bearer ${authToken}` },
      })
      .pipe(
        catchError((error: HttpErrorResponse) => {
          const message = getErrorResponseMessage(error);
          if (message) throw new Error(message);
          throw new Error('Unknown error occured');
        }),
        map((res) => {
          const link = res as Link;
          return link;
        }),
      );
  }
}
