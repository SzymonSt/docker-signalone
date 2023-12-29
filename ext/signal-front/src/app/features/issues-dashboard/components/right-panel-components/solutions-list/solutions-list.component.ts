import { Component, Input } from '@angular/core';
import { IssueDTO } from 'app/shared/interfaces/IssueDTO';

@Component({
  selector: 'app-solutions-list',
  templateUrl: './solutions-list.component.html',
  styleUrls: [ './solutions-list.component.scss' ]
})
export class SolutionsListComponent {
  @Input()
  public activeIssue: IssueDTO;

  constructor() {
  }
}
