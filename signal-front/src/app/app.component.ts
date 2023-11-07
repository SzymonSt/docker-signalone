import { Component } from '@angular/core';
import { ApplicationStateService } from 'app/shared/services/application-state.service';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: [ './app.component.scss' ]
})
export class AppComponent {
  public isLoading$: Observable<boolean>;

  constructor(private applicationStateService: ApplicationStateService) {
    this.isLoading$ = this.applicationStateService.isLoading$;
  }
}
