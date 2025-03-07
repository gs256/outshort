import { version } from '../../package.json';

export const VERSION = version;

export const API_URL = 'http://127.0.0.1:8249';

export const MINUTE = 60;
export const HOUR = MINUTE * 60;
export const DAY = HOUR * 24;
export const WEEK = DAY * 7;
export const MONTH = DAY * 30;
export const YEAR = DAY * 365;

export const LIFETIME_OPTIONS = [
  {
    label: '5 minutes',
    value: MINUTE * 5,
  },
  {
    label: '1 hour',
    value: HOUR,
  },
  {
    label: '5 hours',
    value: HOUR * 5,
  },
  {
    label: '12 hours',
    value: HOUR * 12,
  },
  {
    label: '1 day',
    value: DAY,
  },
  {
    label: '1 week',
    value: WEEK,
  },
  {
    label: '1 month',
    value: MONTH,
  },
  {
    label: '6 month',
    value: MONTH * 6,
  },
  {
    label: '1 year',
    value: YEAR,
  },
  {
    label: 'Never',
    value: 0,
  },
];
