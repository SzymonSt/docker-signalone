import { Component } from '@angular/core';
import { LanguageVersion } from 'app/shared/enum/LanguageVersion';
import { LangugageService } from 'app/shared/services/LanguageService';
import { Observable } from 'rxjs';
import { ApplicationStateService } from 'app/shared/services/ApplicationStateService';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: [ './header.component.scss' ]
})
export class HeaderComponent {
  public LanguageVersion: typeof LanguageVersion = LanguageVersion;
  public activeLanguage$: Observable<LanguageVersion>

  constructor(private languageService: LangugageService, private applicationStateService: ApplicationStateService) {
    this.activeLanguage$ = this.applicationStateService.language$;
  }

  public changeLanguage(language: LanguageVersion): void {
    this.applicationStateService.setLanguage(language);
  }
}
