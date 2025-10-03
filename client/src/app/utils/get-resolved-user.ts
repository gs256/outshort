import { ActivatedRoute } from '@angular/router';
import { User } from '../models/user';
import { toSignal } from '@angular/core/rxjs-interop';

export function getResolvedUser(route: ActivatedRoute): User | undefined {
  const data = toSignal(route.data);
  return data()?.['user'] as User | undefined;
}
