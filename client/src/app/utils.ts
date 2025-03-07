import { HttpErrorResponse } from '@angular/common/http';

export function absoluteRoute(relativeRoute: string): string {
  return `/${relativeRoute}`;
}

export function getErrorResponseMessage(
  error: HttpErrorResponse,
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

export function getShortUrl(alias: string) {
  return `${getOrigin()}/${alias}`;
}

export function getOrigin() {
  return window.location.origin;
}
