import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { catchError, map, Observable } from 'rxjs';
import { API_URL } from '../../constants';

@Injectable({
  providedIn: 'root',
})
export class ApiService {
  private readonly http = inject(HttpClient);

  shorten(url: string): Observable<string> {
    return this.http.post(`${API_URL}/api/v1/shorten`, { url: url }).pipe(
      catchError((error) => {
        throw error;
      }),
      map((res) => {
        if ('alias' in res && typeof res['alias'] == 'string' && res['alias']) {
          return res['alias'] as string;
        }
        throw new Error('Invalid response from server');
      })
    );
  }

  test(): Observable<string> {
    return this.http
      .get(`${API_URL}/api/v1/test`, { responseType: 'text' })
      .pipe(
        catchError((error) => {
          throw error;
        })
      );
  }
}
