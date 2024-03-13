import { CommonModule } from '@angular/common';
import { HTTP_INTERCEPTORS, HttpClientModule } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { AngularSvgIconModule } from 'angular-svg-icon';
import { ErrorHandlerInterceptor } from 'app/shared/interceptors/error-handler.interceptor';
import { LoadingInterceptor } from 'app/shared/interceptors/loading.interceptor';
import {
  ConfigurationDropdownComponent
} from 'app/shared/ui/components/configuration-dropdown/configuration-dropdown.component';
import { ContactPopupComponent } from 'app/shared/ui/components/contact/contact-popup.component';
import { BsDatepickerModule } from 'ngx-bootstrap/datepicker';
import { HeaderComponent } from './ui/components/header/header.component';
import { LoaderComponent } from './ui/components/loader/loader.component';

@NgModule({
  declarations: [
    LoaderComponent,
    HeaderComponent,
    ConfigurationDropdownComponent,
    ContactPopupComponent
  ],
  imports: [
    CommonModule,
    TranslateModule,
    MatProgressSpinnerModule,
    FormsModule,
    MatSlideToggleModule,
    HttpClientModule,
    AngularSvgIconModule.forRoot(),
    MatTooltipModule,
    ReactiveFormsModule,
    BsDatepickerModule

  ],
  exports: [
    LoaderComponent,
    HeaderComponent,
    ConfigurationDropdownComponent,
    ContactPopupComponent
  ],
  providers: [
    { provide: HTTP_INTERCEPTORS, useClass: LoadingInterceptor, multi: true },
    { provide: HTTP_INTERCEPTORS, useClass: ErrorHandlerInterceptor, multi: true }
  ]
})
export class SharedModule {
}
