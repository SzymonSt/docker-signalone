import { Component, Input } from '@angular/core';
import { IssueDTO } from 'app/shared/interfaces/IssueDTO';
import { IssueSeverity } from 'app/shared/enum/IssueSeverity';

@Component({
  selector: 'app-issue-cell',
  templateUrl: './issue-cell.component.html',
  styleUrls: [ './issue-cell.component.scss' ]
})
export class IssueCellComponent {
  public IssueSeverity: typeof IssueSeverity = IssueSeverity;

  @Input()
  public issue: IssueDTO

  constructor() {
  }
}
