import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { SocialAuthService } from '@abacritt/angularx-social-login';
import { AuthStateService } from 'app/auth/services/auth-state.service';
import { ApplicationStateService } from 'app/shared/services/application-state.service';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: [ './login.component.scss' ]
})
export class LoginComponent implements OnInit{
  public loginForm: FormGroup;
  public isSubmitted: boolean = false;
  public activeLocale: string
  constructor(private socialAuthService: SocialAuthService, private authStateService: AuthStateService) {
    this.loginWithGoogle();

  }

  public ngOnInit(): void {
    this.initForm();
  }

  public submitForm(): void {
    this.isSubmitted = true;
    this.loginForm.markAsDirty();
    this.loginForm.markAllAsTouched();
    if (this.loginForm.valid) {
      console.log(this.loginForm.value)
    }
  }

  public loginWithGoogle(): void {
    this.socialAuthService.authState.pipe(takeUntilDestroyed()).subscribe((user) => {
      this.authStateService.loginWithGoogle(user);
    });
  }

  public loginWithGithub(): void {
    this.authStateService.loginWithGithub();
  }


  private initForm(): void {
    this.loginForm = new FormGroup({
      email: new FormControl(null, [Validators.required, Validators.email]),
      password: new FormControl(null, [Validators.required])
    })
  }
}
