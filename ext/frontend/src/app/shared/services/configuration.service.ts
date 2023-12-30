import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from 'environment/environment';

@Injectable({ providedIn: 'root' })
export class ConfigurationService {
  constructor(private httpClient: HttpClient) {
  }

  public changeApiKeyVersion(apiKey: string): Observable<void> {
    return this.httpClient.put<void>(`${environment.apiUrl}/configApiKey`, {apiKey});
  }

}