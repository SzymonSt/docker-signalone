import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { TranslateModule } from '@ngx-translate/core';
import { SharedModule } from './shared/SharedModule';
import { TranslateConfig } from 'app/config/TranslateConfig';
import { HttpClientModule } from '@angular/common/http';

@NgModule({
  declarations: [
    AppComponent,

  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    SharedModule,
    HttpClientModule,
    TranslateModule.forRoot(TranslateConfig)
  ],
  providers: [],
  bootstrap: [ AppComponent ],

})
export class AppModule {
}
