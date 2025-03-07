import { LinkUpsert } from '../models/link-upsert';

export function formToLinkUpsert(formValue: {
  name: string;
  originalUrl: string;
  alias: string | null;
  lifetime: number | null;
}): LinkUpsert {
  return {
    name: formValue.name,
    url: formValue.originalUrl,
    alias: formValue.alias ?? undefined,
    lifetime: formValue.lifetime ?? 0,
  };
}
