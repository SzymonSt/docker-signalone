import { Component, Input } from '@angular/core';
import { IssueDTO } from 'app/shared/interfaces/IssueDTO';

@Component({
  selector: 'app-issues-right-panel',
  templateUrl: './issues-right-panel.component.html',
  styleUrls: [ './issues-right-panel.component.scss' ]
})
export class IssuesRightPanelComponent {
  @Input()
  public activeIssue: IssueDTO;

  constructor() {
  }
}
