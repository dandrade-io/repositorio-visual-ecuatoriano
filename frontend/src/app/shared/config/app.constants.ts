declare global {
  interface Window {
    __env?: { BACKEND_BASE_URL?: string };
  }
}

export const BACKEND_BASE_URL = (
  window.__env?.BACKEND_BASE_URL ?? 'http://localhost:8080'
).replace(/\/$/, '');

export const API_BASE_URL = `${BACKEND_BASE_URL}/api/v1`;
