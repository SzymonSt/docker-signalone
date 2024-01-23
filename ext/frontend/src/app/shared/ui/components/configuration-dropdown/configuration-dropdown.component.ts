import { Component } from '@angular/core';
import { ConfigurationService } from 'app/shared/services/configuration.service';
import { ToastrService } from 'ngx-toastr';
import { TranslateService } from '@ngx-translate/core';

@Component({
  selector: 'app-configuration-dropdown',
  templateUrl: './configuration-dropdown.component.html',
  styleUrls: [ './configuration-dropdown.component.scss' ]
})
export class ConfigurationDropdownComponent {
  public huggingfaceApiKey: string;
  public isSubmitted: boolean;
  constructor(private configurationService: ConfigurationService, private toastrService: ToastrService, private translateService: TranslateService) {
  }

  public submitApiKey(dropdownButton: HTMLButtonElement, dropdownMenu: HTMLElement): void {
    this.isSubmitted = true;
    if (this.huggingfaceApiKey) {
      this.isSubmitted = false;
      this.configurationService.changeApiKeyVersion(this.huggingfaceApiKey).subscribe(() => {
        this.closeDropdownMenu(dropdownButton, dropdownMenu);
        this.toastrService.success(this.translateService.instant('UI.CONFIGURATION_DROPDOWN.SUCCESS'));

      });
    }
  }

  private closeDropdownMenu(dropdownButton: HTMLButtonElement, dropdownMenu: HTMLElement): void {
    dropdownButton.classList.remove('show')
    dropdownButton.ariaExpanded = 'false';
    dropdownMenu.classList.remove('show')
  }

}
