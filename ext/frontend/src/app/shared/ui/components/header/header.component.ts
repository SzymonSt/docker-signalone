import { Component } from '@angular/core';
import { LanguageVersion } from 'app/shared/enum/LanguageVersion';
import { LangugageService } from 'app/shared/services/language.service';
import { Observable } from 'rxjs';
import { ApplicationStateService } from 'app/shared/services/application-state.service';
import { ConfigurationService } from 'app/shared/services/configuration.service';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: [ './header.component.scss' ]
})
export class HeaderComponent {
  public LanguageVersion: typeof LanguageVersion = LanguageVersion;
  public activeLanguage$: Observable<LanguageVersion>

  constructor(private languageService: LangugageService, private applicationStateService: ApplicationStateService, public configurationService: ConfigurationService) {
    this.activeLanguage$ = this.applicationStateService.language$;
  }

  public changeLanguage(language: LanguageVersion): void {
    this.applicationStateService.setLanguage(language);
  }
}
