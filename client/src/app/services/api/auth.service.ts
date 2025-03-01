import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { catchError, map, Observable } from 'rxjs';
import { API_URL } from '../../constants';
import { getErrorResponseMessage } from '../../utils';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private readonly http = inject(HttpClient);

  public signIn(username: string, password: string): Observable<string> {
    return this.http
      .post(`${API_URL}/api/v1/auth/sign-in`, {
        username: username,
        password: password,
      })
      .pipe(
        catchError((error: HttpErrorResponse) => {
          const message = getErrorResponseMessage(error);
          if (message) throw new Error(message);
          throw new Error('Unknown error occured');
        }),
        map((res) => {
          if ('authToken' in res && typeof res['authToken'] === 'string') {
            const authToken = res['authToken'];
            return authToken;
          }
          throw new Error('Unknown error occured');
        }),
      );
  }
}
