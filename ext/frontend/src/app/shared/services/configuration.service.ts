import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable, from, defer } from 'rxjs';
import { environment } from 'environment/environment';
import { AgentStateDTO } from '../interfaces/AgentStateDTO';

@Injectable({ providedIn: 'root' })
export class ConfigurationService {
  constructor() {
  }

  // public getConfiguration(): Observable<AgentStateDTO> {
  //   return from(
  //     this.dockerDesktopClient.extension.vm?.service?.get('/api/control/state')
  //      ?? 
  //      Promise.resolve<AgentStateDTO>({ state: false })
  //     );
  // }
}