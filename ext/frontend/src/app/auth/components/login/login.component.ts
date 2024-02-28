import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { environment } from 'environment/environment.development';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: [ './login.component.scss' ]
})
export class LoginComponent implements OnInit{
  public loginForm: FormGroup;
  public isSubmitted: boolean = false;
  public githubLoginUrl: string = `https://github.com/login/oauth/authorize?client_id=${environment.githubClientId}`;
  public googleLoginUrl: string = `https://accounts.google.com/o/oauth2/v2/auth?scope=openid%20email&nonce=${Math.random() * 100000000}&response_type=id_token&redirect_uri=http://localhost:37001/google-login&client_id=${environment.googleLoginProvider}`;

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

  private initForm(): void {
    this.loginForm = new FormGroup({
      email: new FormControl(null, [Validators.required, Validators.email]),
      password: new FormControl(null, [Validators.required])
    })
  }
}
