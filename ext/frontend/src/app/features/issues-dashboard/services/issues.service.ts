import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from 'environment/environment';
import { IssueSearchCriteriaDTO } from 'app/shared/interfaces/IssueSearchCriteriaDTO';
import { HttpEncoder } from 'app/shared/util/HttpEncoder';
import { NormalizeObjectValue } from 'app/shared/util/NormalizeObjectValue';
import { DetailedIssueDTO } from 'app/shared/interfaces/DetailedIssueDTO';
import { SearchIssuesResponseDTO } from 'app/shared/interfaces/SearchIssuesResponseDTO';
import { RateIssueDTO } from 'app/shared/interfaces/RateIssueDTO';

@Injectable({ providedIn: 'root' })
export class IssuesService {
  constructor(private httpClient: HttpClient) {
  }

  public getIssuesContainers(): Observable<string[]> {
    return this.httpClient.get<string[]>(`${environment.apiUrl}/user/containers`);
  }

  public getIssuesList(searchCriteria?: IssueSearchCriteriaDTO, revokeLoader: boolean = false): Observable<SearchIssuesResponseDTO> {
    if (searchCriteria) {
      if (searchCriteria.startTimestamp) {
        searchCriteria.startTimestamp = new Date(searchCriteria.startTimestamp).toISOString();
      }
      if (searchCriteria.endTimestamp) {
        searchCriteria.endTimestamp = new Date(searchCriteria.endTimestamp).toISOString();
      }

      const params: HttpParams = new HttpParams({
        encoder: new HttpEncoder(),
        fromObject: { ...(NormalizeObjectValue(searchCriteria, [ 'startTimestamp', 'endTimestamp' ]) as any) }
      });

      return this.httpClient.get<SearchIssuesResponseDTO>(`${environment.apiUrl}/user/issues?revokeLoader=${revokeLoader}`, { params });
    } else {
      return this.httpClient.get<SearchIssuesResponseDTO>(`${environment.apiUrl}/user/issues?revokeLoader=${revokeLoader}`);
    }

  }

  public getIssue(issueId: string): Observable<DetailedIssueDTO> {
    return this.httpClient.get<DetailedIssueDTO>(`${environment.apiUrl}/user/issues/${issueId}`);
  }

  public rateIssue(issueId: number, rateIssueData: RateIssueDTO): Observable<void> {
    return this.httpClient.post<void>(`${environment.apiUrl}/user/issues/${issueId}/score`, rateIssueData);
  }

}