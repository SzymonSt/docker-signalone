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
  ],
  providers: [ AlertConfig, BsDatepickerConfig ],
  bootstrap: [ AppComponent ],

})
export class AppModule {
}
