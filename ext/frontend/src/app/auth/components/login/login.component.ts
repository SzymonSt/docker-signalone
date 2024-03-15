import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { ResendVerificationLinkPopupComponent } from 'app/auth/components/resendVerificationLink/resend-verification-link-popup.component';
import { AuthStateService } from 'app/auth/services/auth-state.service';
import { ContactPopupComponent } from 'app/shared/ui/components/contact/contact-popup.component';
import { environment } from 'environment/environment.development';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss'],
})
export class LoginComponent implements OnInit {
  public loginForm: FormGroup;
  public isSubmitted: boolean = false;
  public githubLoginUrl: string = `https://github.com/login/oauth/authorize?client_id=${environment.githubClientId}`;
  public googleLoginUrl: string = `https://accounts.google.com/o/oauth2/v2/auth?scope=openid%20email&nonce=${
    Math.random() * 100000000
  }&response_type=id_token&redirect_uri=http://localhost:37001/google-login&client_id=${
    environment.googleLoginProvider
  }`;

  public constructor(
    private authStateService: AuthStateService,
    private dialog: MatDialog
  ) {}

  public ngOnInit(): void {
    this.initForm();
  }

  public submitForm(): void {
    this.isSubmitted = true;
    this.loginForm.markAsDirty();
    this.loginForm.markAllAsTouched();
    if (this.loginForm.valid) {
      this.authStateService
        .login(
          this.loginForm.get('email').value,
          this.loginForm.get('password').value
        )
        .then()
        .catch(() => this.loginForm.get('password').setValue(null));
    }
  }

  public openResendVerificationLinkModal(): void {
    this.dialog.open(ResendVerificationLinkPopupComponent, {
      width: '500px',
    });
  }

  public openResendVerificationLinkModalKeydown(event: KeyboardEvent): void {
    if (event instanceof KeyboardEvent && event.key === 'Enter') {
      this.dialog.open(ResendVerificationLinkPopupComponent, {
        width: '500px',
      });
    }
  }

  private initForm(): void {
    this.loginForm = new FormGroup({
      email: new FormControl(null, [Validators.required, Validators.email]),
      password: new FormControl(null, [
        Validators.required,
        Validators.minLength(8),
      ]),
    });
  }
}
