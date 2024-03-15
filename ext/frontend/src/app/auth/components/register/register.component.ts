import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { AuthStateService } from 'app/auth/services/auth-state.service';
import { AuthService } from 'app/auth/services/auth.service';
import { Constants } from 'app/config/Constant';
import { confirmPasswordValidator } from 'app/shared/validators/confirm-password.validator';
import { environment } from 'environment/environment.development';
import { ToastrService } from 'ngx-toastr';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: [ './register.component.scss' ]
})
export class RegisterComponent implements OnInit{
  public readonly Constants: typeof Constants = Constants
  public registrationForm: FormGroup;
  public isSubmitted: boolean = false;
  public githubLoginUrl: string = `https://github.com/login/oauth/authorize?client_id=${environment.githubClientId}`;
  public googleLoginUrl: string = `https://accounts.google.com/o/oauth2/v2/auth?scope=openid%20email&nonce=${Math.random() * 100000000}&response_type=id_token&redirect_uri=http://localhost:37001/google-login&client_id=${environment.googleLoginProvider}`;

  public constructor(private authStateService: AuthStateService,
                     private authService: AuthService,
                     private router: Router,
                     private toastrService: ToastrService,
                     private translateService: TranslateService,) {
  }

  public ngOnInit(): void {
    this.initForm();
  }

  public submitForm(): void {
    this.isSubmitted = true;
    this.registrationForm.markAsDirty();
    this.registrationForm.markAllAsTouched();
    console.log(this.registrationForm)
    if (this.registrationForm.valid) {
      this.router.navigateByUrl('/login');
      this.toastrService.success(this.translateService.instant('AUTH.EMAIL_VERIFICATION_LINK_SENT'));
      this.authService.register(this.registrationForm.get('email').value, this.registrationForm.get('password').value).subscribe(() => {
        this.router.navigateByUrl('/login');
        this.toastrService.success(this.translateService.instant('AUTH.EMAIL_VERIFICATION_LINK_SENT'));
      })
    }
  }

  private initForm(): void {
    this.registrationForm = new FormGroup({
      email: new FormControl(null, [Validators.required, Validators.email]),
      password: new FormControl(null, [Validators.required]),
      passwordConfirmation: new FormControl(null, [Validators.required]),
    }, {validators: [confirmPasswordValidator]})
  }
}
