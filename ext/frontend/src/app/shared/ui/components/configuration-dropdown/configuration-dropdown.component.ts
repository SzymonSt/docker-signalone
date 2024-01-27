import { Component} from '@angular/core';
import { ConfigurationService } from 'app/shared/services/configuration.service';
import { ToastrService } from 'ngx-toastr';
import { TranslateService } from '@ngx-translate/core';
import { AgentStateDTO } from 'app/shared/interfaces/AgentStateDTO';
import { MatSlideToggleChange } from '@angular/material/slide-toggle';

@Component({
  selector: 'app-configuration-dropdown',
  templateUrl: './configuration-dropdown.component.html',
  styleUrls: [ './configuration-dropdown.component.scss' ],
})
export class ConfigurationDropdownComponent {
  public agentState: AgentStateDTO = { 
    state: false 
  };
  constructor(private configurationService: ConfigurationService, private toastrService: ToastrService, private translateService: TranslateService) {
  }

  ngOnInit(): void {
    this.configurationService.getCurrentAgentState().subscribe((agentState: AgentStateDTO) => {
      this.agentState = agentState;
    });
  }

  public setAgentState(event: MatSlideToggleChange): void {
    this.agentState.state = event.checked
    console.log(this.agentState);
    this.configurationService.setAgentState(this.agentState).subscribe(() => {
      this.toastrService.success(this.translateService.instant('configuration.agentStateUpdated'));
    });
  }

  private closeDropdownMenu(dropdownButton: HTMLButtonElement, dropdownMenu: HTMLElement): void {
    dropdownButton.classList.remove('show')
    dropdownButton.ariaExpanded = 'false';
    dropdownMenu.classList.remove('show')
  }

}
