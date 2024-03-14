import { Clipboard } from '@angular/cdk/clipboard';
import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormControl, FormGroup, Validators } from '@angular/forms';
import { MatDialogRef } from '@angular/material/dialog';
import { TranslateService } from '@ngx-translate/core';
import { AuthService } from 'app/auth/services/auth.service';
import { ToastrService } from 'ngx-toastr';

@Component({
  selector: 'app-resend-verification-link-popup',
  templateUrl: './resend-verification-link-popup.component.html',
  styleUrls: [ './resend-verification-link-popup.component.scss' ]
})
export class ResendVerificationLinkPopupComponent implements OnInit {
  public resendVerificationLinkForm: FormGroup;
  public isSubmitted: boolean = false;

  public get emailControl(): AbstractControl {
    return this.resendVerificationLinkForm.get('email');
  }

  constructor(public dialogRef: MatDialogRef<ResendVerificationLinkPopupComponent>,
              private clipboard: Clipboard,
              private toastrService: ToastrService,
              private translateService: TranslateService,
              private authService: AuthService) {
  }

  public ngOnInit(): void {
    this.initForm();
  }

  public submitContact(): void {
    this.isSubmitted = true;
    this.resendVerificationLinkForm.markAsDirty();
    this.resendVerificationLinkForm.markAllAsTouched();
    if (this.resendVerificationLinkForm.valid) {
      this.authService.resendVerificationLink(this.resendVerificationLinkForm.value.email).subscribe(res => {
        this.toastrService.success(this.translateService.instant('AUTH.VERIFICATION_LINK_RESEND_SUCCESS'));
        this.close();
      })
    }
  }

  public close(): void {
    this.dialogRef.close();
  }

  private initForm(): void {
    this.resendVerificationLinkForm = new FormGroup({
      email: new FormControl(null, [Validators.email, Validators.required])
    });
  }
}
