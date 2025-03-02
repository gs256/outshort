const AUTH_TOKEN_KEY = 'authToken';

export const storage = {
  get authToken(): string {
    return localStorage.getItem(AUTH_TOKEN_KEY)?.trim() ?? '';
  },

  set authToken(value: string) {
    localStorage.setItem(AUTH_TOKEN_KEY, value.trim());
  },
};
