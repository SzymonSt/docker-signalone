import { NgModule } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { CommonModule } from '@angular/common';
import { LoaderComponent } from './ui/components/loader/loader.component';
import { HeaderComponent } from './ui/components/header/header.component';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { LoadingInterceptor } from 'app/shared/interceptors/loading.interceptor';
import { HTTP_INTERCEPTORS } from '@angular/common/http';
import { ErrorHandlerInterceptor } from 'app/shared/interceptors/error-handler.interceptor';

@NgModule({
  declarations: [
    LoaderComponent,
    HeaderComponent
  ],
  imports: [
    CommonModule,
    TranslateModule,
    MatProgressSpinnerModule
  ],
  exports: [
    LoaderComponent,
    HeaderComponent
  ],
  providers: [
    { provide: HTTP_INTERCEPTORS, useClass: LoadingInterceptor, multi: true },
    { provide: HTTP_INTERCEPTORS, useClass: ErrorHandlerInterceptor, multi: true }
  ]
})
export class SharedModule {
}
