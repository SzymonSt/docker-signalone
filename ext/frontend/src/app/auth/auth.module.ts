import { NgModule } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { CommonModule } from '@angular/common';
import { GoogleLoginComponent } from 'app/auth/components/googleLogin/google-login.component';
import { SharedModule } from 'app/shared/SharedModule';
import { NgSelectModule } from '@ng-select/ng-select';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { BsDatepickerModule } from 'ngx-bootstrap/datepicker';
import { AngularSvgIconModule } from 'angular-svg-icon';
import { AuthRoutingModule } from 'app/auth/auth-routing.module';
import { LoginComponent } from 'app/auth/components/login/login.component';
import { HTTP_INTERCEPTORS } from '@angular/common/http';
import { AuthInterceptor } from 'app/shared/interceptors/auth.interceptor';
import { GoogleSigninButtonModule } from '@abacritt/angularx-social-login';
import { GithubLoginComponent } from 'app/auth/components/githubLogin/github-login.component';

@NgModule({
  declarations: [ LoginComponent, GithubLoginComponent, GoogleLoginComponent ],
  imports: [
    CommonModule,
    TranslateModule,
    SharedModule,
    NgSelectModule,
    FormsModule,
    GoogleSigninButtonModule,
    ReactiveFormsModule,
    NgbModule,
    BsDatepickerModule,
    AuthRoutingModule,
    AngularSvgIconModule.forRoot()
  ],
  exports: [],
  providers: [
    { provide: HTTP_INTERCEPTORS, useClass: AuthInterceptor, multi: true }
  ]
})
export class AuthModule {
}
