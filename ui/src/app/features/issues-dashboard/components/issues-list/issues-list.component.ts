import { Component, Input } from '@angular/core';
import { IssueDTO } from 'app/shared/interfaces/IssueDTO';

@Component({
  selector: 'app-issues-list',
  templateUrl: './issues-list.component.html',
  styleUrls: [ './issues-list.component.scss' ]
})
export class IssuesListComponent {
  @Input()
  public issues: IssueDTO[];

  constructor() {
  }
}
  