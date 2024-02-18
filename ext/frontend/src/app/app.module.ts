import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { TranslateModule } from '@ngx-translate/core';
import { SharedModule } from './shared/SharedModule';
import { TranslateConfig } from 'app/config/TranslateConfig';
import { HttpClientModule } from '@angular/common/http';
import { BsDatepickerConfig, BsDatepickerModule } from 'ngx-bootstrap/datepicker';
import { AlertConfig } from 'ngx-bootstrap/alert';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import * as moment from 'moment';
import { ToastrModule } from 'ngx-toastr';
import { AuthModule } from 'app/auth/auth.module';
import { GoogleLoginProvider, SocialAuthServiceConfig } from '@abacritt/angularx-social-login';

@NgModule({
  declarations: [
    AppComponent,

  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    SharedModule,
    HttpClientModule,
    TranslateModule.forRoot(TranslateConfig),
    BrowserAnimationsModule,
    BsDatepickerModule.forRoot(),
    ToastrModule.forRoot({
      positionClass: 'toast-bottom-center',
      preventDuplicates: true,
      extendedTimeOut: moment.duration(3, 'seconds').as('milliseconds'),
      enableHtml: true
    }),
    AuthModule,
  ],
  providers: [ AlertConfig, BsDatepickerConfig,
    {
      provide: 'SocialAuthServiceConfig',
      useValue: {
        autoLogin: false,
        providers: [
          {
            id: GoogleLoginProvider.PROVIDER_ID,
            provider: new GoogleLoginProvider('359898712853-rcec16l3ivs24rod4hb2kcsr9qvotf3t.apps.googleusercontent.com'),
          },
        ],
      } as SocialAuthServiceConfig,
    },],
  bootstrap: [ AppComponent ],

})
export class AppModule {
}
