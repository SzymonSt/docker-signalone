import { Component } from '@angular/core';
import { ApplicationStateService } from 'app/shared/services/application-state.service';
import { Observable } from 'rxjs';
import { ConfigurationService } from 'app/shared/services/configuration.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: [ './app.component.scss' ]
})
export class AppComponent {
  public isLoading$: Observable<boolean>;

  constructor(private applicationStateService: ApplicationStateService, private configurationService: ConfigurationService) {
    this.isLoading$ = this.applicationStateService.isLoading$;
    this.configurationService.getCurrentAgentState();
  }
}
