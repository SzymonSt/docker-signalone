import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from 'environment/environment';
import { IssueDTO } from 'app/shared/interfaces/IssueDTO';
import { IssueSearchCriteriaDTO } from 'app/shared/interfaces/IssueSearchCriteriaDTO';
import { HttpEncoder } from 'app/shared/util/HttpEncoder';
import { NormalizeObjectValue } from 'app/shared/util/NormalizeObjectValue';

@Injectable({ providedIn: 'root' })
export class IssuesService {
  constructor(private httpClient: HttpClient) {
  }

  public getIssuesContainers(): Observable<string[]> {
    return this.httpClient.get<string[]>(`${environment.apiUrl}/containers`);
  }

  public getIssuesList(searchCriteria?: IssueSearchCriteriaDTO): Observable<IssueDTO[]> {
    if (searchCriteria) {
      const params: HttpParams = new HttpParams({
        encoder: new HttpEncoder(),
        fromObject: { ...(NormalizeObjectValue(searchCriteria) as any) }
      });

      return this.httpClient.get<IssueDTO[]>(`${environment.apiUrl}/issues`, { params });
    } else {
      return this.httpClient.get<IssueDTO[]>(`${environment.apiUrl}/issues`);
    }

  }

}