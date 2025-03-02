import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { catchError, map, Observable, throwError } from 'rxjs';
import { API_URL } from '../../constants';
import { getErrorResponseMessage } from '../../utils';
import { User } from '../../models/user';
import { storage } from '../../storage';
import { Auth } from '../../models/auth';

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
          const auth = res as Auth;
          storage.authToken = auth.authToken;
          return auth.authToken;
        }),
      );
  }

  public getUserInfo(): Observable<User> {
    const authToken = storage.authToken;
    if (authToken.length === 0) {
      return throwError(() => new Error('No auth token'));
    }
    return this.http
      .get(`${API_URL}/api/v1/auth/user-info`, {
        headers: { Authorization: `Bearer ${authToken}` },
      })
      .pipe(
        catchError((error: HttpErrorResponse) => {
          const message = getErrorResponseMessage(error);
          if (message) throw new Error(message);
          throw new Error('Unknown error occured');
        }),
        map((res) => {
          const user = res as User;
          return user;
        }),
      );
  }
}
