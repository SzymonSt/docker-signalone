import { Component, OnChanges, OnInit, SimpleChanges } from '@angular/core';
import { ConfigurationService } from 'app/shared/services/configuration.service';
import { AgentStateDTO } from 'app/shared/interfaces/AgentStateDTO';
import { MatSlideToggleChange } from '@angular/material/slide-toggle';

@Component({
  selector: 'app-configuration-dropdown',
  templateUrl: './configuration-dropdown.component.html',
  styleUrls: [ './configuration-dropdown.component.scss' ]
})
export class ConfigurationDropdownComponent implements OnInit {
  public agentState: AgentStateDTO;
  constructor(private configurationService: ConfigurationService) {
  }

  public ngOnInit(): void {
    this.agentState = new AgentStateDTO(this.configurationService.currentAgentState);
  }

  public setAgentState(): void {
    console.log(this.agentState)
    this.configurationService.setAgentState(this.agentState);
  }

}
