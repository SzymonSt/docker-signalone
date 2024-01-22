import { Component, Input } from '@angular/core';
import { IssueDTO } from 'app/shared/interfaces/IssueDTO';
import { DetailedIssueDTO } from 'app/shared/interfaces/DetailedIssueDTO';

@Component({
  selector: 'app-solutions-list',
  templateUrl: './solutions-list.component.html',
  styleUrls: [ './solutions-list.component.scss' ]
})
export class SolutionsListComponent {
  @Input()
  public activeIssue: DetailedIssueDTO;

  constructor() {
  }

  public goToSolutionSource(url : string): void {
    window.open(url, "_blank")
  }
}
