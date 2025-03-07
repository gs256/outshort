export interface LinkUpsert {
  name: string;
  url: string;
  alias?: string;
  lifetime: number;
}
