import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from 'environment/environment';
import { AgentStateDTO } from '../interfaces/AgentStateDTO';

@Injectable({ providedIn: 'root' })
export class ConfigurationService {
  constructor(private httpClient: HttpClient) {
  }

  public getCurrentAgentState(): Observable<AgentStateDTO> {
    return this.httpClient.get<AgentStateDTO>(`${environment.agentApiUrl}/control/state`);
  }
  
  public setAgentState(agentStatePayload: AgentStateDTO): Observable<void> {
    return this.httpClient.post<void>(`${environment.agentApiUrl}/control/state`, agentStatePayload);
  }
}