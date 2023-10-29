import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from 'environment/environment';

@Injectable({ providedIn: 'root' })
export class IssuesService {
  constructor(private httpClient: HttpClient) {
  }

  public getIssuesContainers(): Observable<string[]> {
    console.log('TEST')
    return this.httpClient.get<string[]>(`${environment.apiUrl}/containers`);
  }

}