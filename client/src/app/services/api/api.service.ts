import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { catchError, map, Observable } from 'rxjs';
import { API_URL } from '../../constants';

@Injectable({
  providedIn: 'root',
})
export class ApiService {
  private readonly http = inject(HttpClient);

  public shorten(url: string): Observable<string> {
    return this.http.post(`${API_URL}/api/v1/shorten`, { url: url }).pipe(
      catchError((error: HttpErrorResponse) => {
        const message = this.getErrorResponseMessage(error);
        if (message) throw new Error(message);
        throw new Error('Unknown error occured');
      }),
      map((res) => {
        if ('alias' in res && typeof res['alias'] == 'string' && res['alias']) {
          return res['alias'] as string;
        }
        throw new Error('Invalid response from server');
      })
    );
  }

  public test(): Observable<string> {
    return this.http
      .get(`${API_URL}/api/v1/test`, { responseType: 'text' })
      .pipe(
        catchError((error) => {
          throw error;
        })
      );
  }

  private getErrorResponseMessage(
    error: HttpErrorResponse
  ): string | undefined {
    if (error.status == 0) {
      return undefined;
    }
    const errorBody = error.error;
    if ('error' in errorBody && typeof errorBody['error'] == 'string') {
      return errorBody['error'] as string;
    }
    return undefined;
  }
}
