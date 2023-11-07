import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { LangugageService } from './language.service';
import { LanguageVersion } from 'app/shared/enum/LanguageVersion';
import { ApplicationConfig } from 'app/config/ApplicationConfig';

@Injectable({ providedIn: 'root' })
export class ApplicationStateService {
  public isLoading$: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public language$: BehaviorSubject<LanguageVersion> = new BehaviorSubject<LanguageVersion>(ApplicationConfig.defaultLanguage as LanguageVersion);

  public constructor(private languageService: LangugageService) {
  }

  public get isLoading(): boolean {
    return this.isLoading$.value;
  }

  public set isLoading(value: boolean) {
    this.isLoading$.next(value);
  }

  public get language(): LanguageVersion {
    return this.language$.value;
  }

  public set language(value: LanguageVersion) {
    this.language$.next(value);
  }

  public setLanguage(language: LanguageVersion): void {
    if (!this.language) {
      this.languageService.initialize();
    }

    this.languageService.setLanguage(language);
    this.language = language;
  }

}