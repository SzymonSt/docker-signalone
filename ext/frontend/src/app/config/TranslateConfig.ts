import { TranslateLoader, TranslateModuleConfig } from '@ngx-translate/core';
import { HttpClient } from '@angular/common/http';
import { from, Observable } from 'rxjs';

// newer, webpack approach (compiled-in direct import (webpackMode: 'eager') or lazy import (webpackMode: 'lazy') + cache busting during production build by webpack)
const WebpackTranslateLoaderFactory: () => TranslateLoader = () => {
  class WebpackTranslateLoader implements TranslateLoader {
    public getTranslation(lang: string): Observable<any> {
      return from(import(
        /* webpackChunkName: "[request]" */
        /* webpackMode: "eager" */
        /* webpackPrefetch: true */
        /* webpackPreload: true */
        `../../assets/locale/${lang}.json`
        ));
    }
  }

  return new WebpackTranslateLoader();
};

export const TranslateConfig: TranslateModuleConfig = {
  loader: {
    provide: TranslateLoader,
    useFactory: WebpackTranslateLoaderFactory,
    deps: [ HttpClient ]
  },
  defaultLanguage: 'en'
};
