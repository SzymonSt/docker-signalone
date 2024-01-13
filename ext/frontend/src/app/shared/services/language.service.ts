import { Injectable } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { registerLocaleData } from '@angular/common';
import localePL from '@angular/common/locales/pl';
import localeEN from '@angular/common/locales/en-GB';
import { defineLocale } from 'ngx-bootstrap/chronos';
import { enGbLocale as localeEN_ngx_bootstrap, plLocale as localePL_ngx_bootstrap } from 'ngx-bootstrap/locale';
import { BsLocaleService } from 'ngx-bootstrap/datepicker';
import { ApplicationConfig } from 'app/config/ApplicationConfig';
import { LanguageVersion } from 'app/shared/enum/LanguageVersion';
import * as moment from 'moment';

@Injectable({ providedIn: 'root' })
export class LangugageService {

  public constructor(private translateService: TranslateService, private localeService: BsLocaleService) {
    this.setLanguage(ApplicationConfig.defaultLanguage as LanguageVersion);
  }

  public initialize(): void {
    this.translateService.addLangs([ LanguageVersion.EN, LanguageVersion.PL ]);
  }

  public setLanguage(language: LanguageVersion): void {
    switch (language) {
      case LanguageVersion.EN: {
        this.translateService.use(LanguageVersion.EN);
        registerLocaleData(localeEN, 'en-GB');
        defineLocale('en-gb', localeEN_ngx_bootstrap);
        this.localeService.use('en-gb');
        moment.locale('en-gb');
        break;
      }
      case LanguageVersion.PL: {
        this.translateService.use(LanguageVersion.PL);
        registerLocaleData(localePL, 'pl');
        defineLocale('pl', localePL_ngx_bootstrap);
        this.localeService.use('pl');
        moment.locale('pl');
        break;
      }
    }
  }
}