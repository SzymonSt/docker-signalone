import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from 'environment/environment';
import { IssueDTO } from 'app/shared/interfaces/IssueDTO';
import { IssueSearchCriteriaDTO } from 'app/shared/interfaces/IssueSearchCriteriaDTO';
import { HttpEncoder } from 'app/shared/util/HttpEncoder';
import { NormalizeObjectValue } from 'app/shared/util/NormalizeObjectValue';
import { DetailedIssueDTO } from 'app/shared/interfaces/DetailedIssueDTO';
import { SearchIssuesResponseDTO } from 'app/shared/interfaces/SearchIssuesResponseDTO';

@Injectable({ providedIn: 'root' })
export class IssuesService {
  constructor(private httpClient: HttpClient) {
  }

  public getIssuesContainers(): Observable<string[]> {
    return this.httpClient.get<string[]>(`${environment.apiUrl}/containers`);
  }

  public getIssuesList(searchCriteria?: IssueSearchCriteriaDTO): Observable<SearchIssuesResponseDTO> {
    if (searchCriteria) {
      if (searchCriteria.startTimestamp) {
        searchCriteria.startTimestamp = new Date(searchCriteria.startTimestamp).toISOString().split('T')[0];
      }
      if (searchCriteria.endTimestamp) {
        searchCriteria.endTimestamp = new Date(searchCriteria.endTimestamp).toISOString().split('T')[0];
      }

      const params: HttpParams = new HttpParams({
        encoder: new HttpEncoder(),
        fromObject: { ...(NormalizeObjectValue(searchCriteria, [ 'startTimestamp', 'endTimestamp' ]) as any) }
      });

      return this.httpClient.get<SearchIssuesResponseDTO>(`${environment.apiUrl}/user/issues`, { params });
    } else {
      return this.httpClient.get<SearchIssuesResponseDTO>(`${environment.apiUrl}/user/issues`);
    }

  }

  public getIssue(issueId: string): Observable<DetailedIssueDTO> {
    return this.httpClient.get<DetailedIssueDTO>(`${environment.apiUrl}/user/issues/${issueId}`);
  }

}