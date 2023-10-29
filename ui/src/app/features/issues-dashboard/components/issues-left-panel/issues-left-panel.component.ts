import { Component, EventEmitter, Input, Output } from '@angular/core';
import { IssueSearchCriteriaDTO } from 'app/shared/interfaces/IssueSearchCriteriaDTO';

@Component({
  selector: 'app-issues-left-panel',
  templateUrl: './issues-left-panel.component.html',
  styleUrls: [ './issues-left-panel.component.scss' ]
})
export class IssuesLeftPanelComponent {
  @Input()
  public containers: string[];
  @Output()
  public criteriaChanged: EventEmitter<IssueSearchCriteriaDTO> = new EventEmitter<IssueSearchCriteriaDTO>();

  constructor() {
  }
}
