import { Component, EventEmitter, Input, Output } from '@angular/core';
import { IssueDTO } from 'app/shared/interfaces/IssueDTO';

@Component({
  selector: 'app-issues-list',
  templateUrl: './issues-list.component.html',
  styleUrls: ['./issues-list.component.scss'],
})
export class IssuesListComponent {
  @Input()
  public issues: IssueDTO[];

  @Output()
  public viewIssue: EventEmitter<IssueDTO> = new EventEmitter<IssueDTO>();

  public selectedIssueId: string;

  constructor() {}

  public selectIssue(issue: IssueDTO): void {
    this.selectedIssueId = issue.id;
    this.viewIssue.emit(issue);
  }

  public selectIssueKeydown(issue: IssueDTO, event: KeyboardEvent): void {
    if (
      event instanceof KeyboardEvent &&
      (event.key === 'Enter' || event.key === ' ')
    ) {
      if (event.key === ' ') event.preventDefault();

      this.selectedIssueId = issue.id;
      this.viewIssue.emit(issue);
    }
  }
}
