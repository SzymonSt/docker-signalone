import { SocialAuthService } from '@abacritt/angularx-social-login';
import { Component, OnInit } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { AuthStateService } from 'app/auth/services/auth-state.service';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: [ './login.component.scss' ]
})
export class LoginComponent implements OnInit{
  public loginForm: FormGroup;
  public isSubmitted: boolean = false;
  public activeLocale: string
  public githubLoginUrl: string = `https://github.com/login/oauth/authorize?client_id=6c88a4f9d4868879974e`;
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

  private initForm(): void {
    this.loginForm = new FormGroup({
      email: new FormControl(null, [Validators.required, Validators.email]),
      password: new FormControl(null, [Validators.required])
    })
  }
}
