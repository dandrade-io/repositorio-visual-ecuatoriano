import { HttpBackend, HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable, tap } from 'rxjs';

export type Env = Record<string, string>;

@Injectable({ providedIn: 'root' })
export class EnvironmentConfigService {
  private http = new HttpClient(inject(HttpBackend));
  private env: Record<string, string> = {};

  static readonly KC_URL = 'KC_URL';
  static readonly KC_REALM = 'KC_REALM';
  static readonly KC_CLIENT_ID = 'KC_CLIENT_ID';
  static readonly API_URL = 'API_URL';

  load(): Observable<Env> {
    return this.http.get<Env>('/assets/env.json').pipe(
      tap((data) => {
        this.env = data;
      }),
    );
  }

  get(key: string): string {
    return this.env?.[key];
  }
}
