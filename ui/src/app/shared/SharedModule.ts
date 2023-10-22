import { NgModule } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { CommonModule } from '@angular/common';
import { LoaderComponent } from './ui/components/loader/loader.component';
import { HeaderComponent } from './ui/components/header/header.component';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';

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
  ]
})
export class SharedModule {
}
