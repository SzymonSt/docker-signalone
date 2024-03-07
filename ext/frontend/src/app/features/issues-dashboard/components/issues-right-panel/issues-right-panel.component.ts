import { Component, EventEmitter, Input, Output } from '@angular/core';
import { DetailedIssueDTO, DetailedIssueScore } from 'app/shared/interfaces/DetailedIssueDTO';

@Component({
  selector: 'app-issues-right-panel',
  templateUrl: './issues-right-panel.component.html',
  styleUrls: [ './issues-right-panel.component.scss' ]
})
export class IssuesRightPanelComponent {
  @Input()
  public activeIssue: DetailedIssueDTO;
  @Output()
  public scoreSelected: EventEmitter<DetailedIssueScore> = new EventEmitter<DetailedIssueScore>();
  @Output()
  public markIssueAsResolved: EventEmitter<void> = new EventEmitter<void>();
  @Output()
  public regenerateIssue: EventEmitter<void> = new EventEmitter<void>();
  constructor() {
  }
}
