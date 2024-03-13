import { Clipboard } from '@angular/cdk/clipboard';
import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormControl, FormGroup, Validators } from '@angular/forms';
import { MatDialogRef } from '@angular/material/dialog';
import { TranslateService } from '@ngx-translate/core';
import { ContactRequestDTO } from 'app/shared/interfaces/ContactRequestDTO';
import { ContactService } from 'app/shared/services/contact.service';
import { ToastrService } from 'ngx-toastr';

@Component({
  selector: 'app-contact-popup',
  templateUrl: './contact-popup.component.html',
  styleUrls: [ './contact-popup.component.scss' ]
})
export class ContactPopupComponent implements OnInit {
  public contactForm: FormGroup;
  public isSubmitted: boolean = false;

  public get emailControl(): AbstractControl {
    return this.contactForm.get('email');
  }

  public get messageTitleControl(): AbstractControl {
    return this.contactForm.get('messageTitle');
  }

  public get messageContentControl(): AbstractControl {
    return this.contactForm.get('messageContent');
  }

  constructor(public dialogRef: MatDialogRef<ContactPopupComponent>,
              private clipboard: Clipboard,
              private toastrService: ToastrService,
              private translateService: TranslateService,
              private contactService: ContactService) {
  }

  public ngOnInit(): void {
    this.initForm();
  }

  public submitContact(): void {
    this.isSubmitted = true;
    this.contactForm.markAsDirty();
    this.contactForm.markAllAsTouched();
    if (this.contactForm.valid) {
      this.contactService.sendContactMessage(new ContactRequestDTO(
        this.contactForm.value.email,
        this.contactForm.value.messageTitle,
        this.contactForm.value.messageContent
      )).subscribe(res => {
        this.toastrService.success(this.translateService.instant('UI.CONTACT_POPUP.MESSAGE_SEND_SUCCESS'));
        this.close();
      })
    }
  }

  public close(): void {
    this.dialogRef.close();
  }

  public copyDiscordUrl(): void {
    this.clipboard.copy('https://discord.gg/vAZrxKs5f6');
    this.toastrService.success(this.translateService.instant('UI.CONTACT_POPUP.COPY_DISCORD_URL_SUCCESS'));
  }

  public copyGithubUrl(): void {
    this.clipboard.copy('https://github.com/Signal0ne');
    this.toastrService.success(this.translateService.instant('UI.CONTACT_POPUP.COPY_GITHUB_URL_SUCCESS'));
  }

  public copyLinkedinUrl(): void {
    this.clipboard.copy('https://www.linkedin.com/company/signal0ne/');
    this.toastrService.success(this.translateService.instant('UI.CONTACT_POPUP.COPY_LINKEDIN_URL_SUCCESS'));
  }

  private initForm(): void {
    this.contactForm = new FormGroup({
      email: new FormControl(null, [Validators.email, Validators.required]),
      messageTitle: new FormControl(null, [Validators.required]),
      messageContent: new FormControl(null, [Validators.required]),
    });
  }
}
