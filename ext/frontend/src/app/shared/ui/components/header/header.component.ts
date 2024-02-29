import { Component } from '@angular/core';
import { LanguageVersion } from 'app/shared/enum/LanguageVersion';
import { LangugageService } from 'app/shared/services/language.service';
import { Observable } from 'rxjs';
import { ApplicationStateService } from 'app/shared/services/application-state.service';
import { ConfigurationService } from 'app/shared/services/configuration.service';
import { AuthStateService } from 'app/auth/services/auth-state.service';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss'],
})
export class HeaderComponent {
  public LanguageVersion: typeof LanguageVersion = LanguageVersion;
  public activeLanguage$: Observable<LanguageVersion>;
  public isLoggedIn$: Observable<boolean>;

  constructor(
    private languageService: LangugageService,
    private applicationStateService: ApplicationStateService,
    protected configurationService: ConfigurationService,
    private authStateService: AuthStateService
  ) {
    this.activeLanguage$ = this.applicationStateService.language$;
    this.isLoggedIn$ = this.authStateService.isLoggedIn$;
  }

  public changeLanguage(language: LanguageVersion): void {
    this.applicationStateService.setLanguage(language);
  }

  public changeLanguageKeydown(
    language: LanguageVersion,
    event: KeyboardEvent
  ): void {
    if (
      event instanceof KeyboardEvent &&
      (event.key === 'Enter' || event.key === ' ')
    ) {
      if (event.key === ' ') event.preventDefault();

      this.applicationStateService.setLanguage(language);
    }
  }

  public logOut(): void {
    this.authStateService.logout();
  }
}
