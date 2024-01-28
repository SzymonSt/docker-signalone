import { NgModule } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { CommonModule } from '@angular/common';
import { MatSlideToggleModule } from '@angular/material/slide-toggle'; 
import { LoaderComponent } from './ui/components/loader/loader.component';
import { HeaderComponent } from './ui/components/header/header.component';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { LoadingInterceptor } from 'app/shared/interceptors/loading.interceptor';
import { HTTP_INTERCEPTORS, HttpClientModule } from '@angular/common/http';
import { ErrorHandlerInterceptor } from 'app/shared/interceptors/error-handler.interceptor';
import {
  ConfigurationDropdownComponent
} from 'app/shared/ui/components/configuration-dropdown/configuration-dropdown.component';
import { FormsModule } from '@angular/forms';
import { AngularSvgIconModule, provideAngularSvgIcon } from 'angular-svg-icon';

@NgModule({
  declarations: [
    LoaderComponent,
    HeaderComponent,
    ConfigurationDropdownComponent
  ],
  imports: [
    CommonModule,
    TranslateModule,
    MatProgressSpinnerModule,
    FormsModule,
    MatSlideToggleModule,
    HttpClientModule,
    AngularSvgIconModule.forRoot()
  ],
  exports: [
    LoaderComponent,
    HeaderComponent,
    ConfigurationDropdownComponent
  ],
  providers: [
    { provide: HTTP_INTERCEPTORS, useClass: LoadingInterceptor, multi: true },
    { provide: HTTP_INTERCEPTORS, useClass: ErrorHandlerInterceptor, multi: true }
  ]
})
export class SharedModule {
}
