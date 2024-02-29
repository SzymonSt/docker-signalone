import { HttpClientModule } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { TranslateModule } from '@ngx-translate/core';
import { AuthModule } from 'app/auth/auth.module';
import { TranslateConfig } from 'app/config/TranslateConfig';
import * as moment from 'moment';
import { AlertConfig } from 'ngx-bootstrap/alert';
import { BsDatepickerConfig, BsDatepickerModule } from 'ngx-bootstrap/datepicker';
import { MarkdownModule } from 'ngx-markdown';
import { ToastrModule } from 'ngx-toastr';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { SharedModule } from './shared/SharedModule';

// @ts-ignore
// @ts-ignore
// @ts-ignore
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
    AuthModule
  ],
  providers: [ AlertConfig, BsDatepickerConfig],
  bootstrap: [ AppComponent ],

})
export class AppModule {
}
