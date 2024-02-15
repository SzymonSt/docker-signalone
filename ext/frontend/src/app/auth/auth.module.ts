import { NgModule } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { CommonModule } from '@angular/common';
import { SharedModule } from 'app/shared/SharedModule';
import { NgSelectModule } from '@ng-select/ng-select';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { BsDatepickerModule } from 'ngx-bootstrap/datepicker';
import { AngularSvgIconModule } from 'angular-svg-icon';
import { AuthRoutingModule } from 'app/auth/auth-routing.module';
import { LoginComponent } from 'app/auth/components/login/login.component';

@NgModule({
  declarations: [ LoginComponent ],
  imports: [
    CommonModule,
    TranslateModule,
    SharedModule,
    NgSelectModule,
    FormsModule,
    ReactiveFormsModule,
    NgbModule,
    BsDatepickerModule,
    AuthRoutingModule,
    AngularSvgIconModule.forRoot()
  ],
  exports: []
})
export class AuthModule {
}
