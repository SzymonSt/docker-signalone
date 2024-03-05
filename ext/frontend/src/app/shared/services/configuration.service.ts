import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, Observable, tap } from 'rxjs';
import { environment } from 'environment/environment';
import { AgentStateDTO } from '../interfaces/AgentStateDTO';
import { ToastrService } from 'ngx-toastr';
import { TranslateService } from '@ngx-translate/core';
import { agentAuthDataDTO } from '../interfaces/AgentAuthDataDTO';

@Injectable({ providedIn: 'root' })
export class ConfigurationService {
  public currentAgentState$: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public isCurrentAgentStateInitialized: boolean = false;
  
  public get currentAgentState(): boolean {
    return  this.currentAgentState$.value;
  }

  public set currentAgentState(value: boolean) {
    this.currentAgentState$.next(value);
  }

  constructor(private httpClient: HttpClient, private toastrService: ToastrService, private translateService: TranslateService) {
  }

  public getCurrentAgentState(): void {
    this.httpClient.get<AgentStateDTO>(`${environment.agentApiUrl}/control/state`).subscribe((agentState: AgentStateDTO) => {
      this.currentAgentState = agentState.state;
      this.isCurrentAgentStateInitialized = true;
    })
  }
  
  public setAgentState(agentStatePayload: AgentStateDTO): void {
    this.httpClient.post<void>(`${environment.agentApiUrl}/control/state`, agentStatePayload).subscribe(() => {
      this.currentAgentState = agentStatePayload.state;
      if (this.currentAgentState) {
        this.toastrService.success(this.translateService.instant('CONFIGURATION.AGENT_STATE_ACTIVATED'));
      } else {
        this.toastrService.success(this.translateService.instant('CONFIGURATION.AGENT_STATE_DEACTIVATED'));
      }

    });
  }

  public setAgentAuthData(agentAuthData: agentAuthDataDTO): void {
    this.httpClient.post<void>(`${environment.agentApiUrl}/control/auth_data`, agentAuthData).subscribe(() => {
      this.toastrService.success(this.translateService.instant('CONFIGURATION.AGENT_AUTH_DATA_UPDATED'));
    });
  }
}