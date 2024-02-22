import { Component, OnInit } from '@angular/core';
import { ApplicationStateService } from 'app/shared/services/application-state.service';
import { Observable } from 'rxjs';
import { ConfigurationService } from 'app/shared/services/configuration.service';
import { AuthStateService } from 'app/auth/services/auth-state.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: [ './app.component.scss' ]
})
export class AppComponent implements OnInit{
  public isLoading$: Observable<boolean>;

  constructor(private applicationStateService: ApplicationStateService, private configurationService: ConfigurationService, private authStateService: AuthStateService, private router: Router) {
    this.isLoading$ = this.applicationStateService.isLoading$;
  }

  public ngOnInit(): void {
    this.authStateService.recoverToken().then(() => {
      this.configurationService.getCurrentAgentState();
      this.router.navigateByUrl('/issues-dashboard')
    }).catch(err => {
      if (!this.router.url.includes('login')) {
        this.router.navigateByUrl('/login')
      }
    })
  }
}
