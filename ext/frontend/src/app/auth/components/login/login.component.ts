import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { Constants } from 'app/config/Constant';
import { dateRangeValidator } from 'app/shared/validators/date-range.validator';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: [ './login.component.scss' ]
})
export class LoginComponent implements OnInit{
  public loginForm: FormGroup;
  public isSubmitted = false;
  constructor() {
  }

  public ngOnInit(): void {
    this.initForm();
  }

  public submitForm(): void {
    this.isSubmitted = true;
    this.loginForm.markAsDirty();
    this.loginForm.markAllAsTouched();
    console.log('TEST', this.loginForm)
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
