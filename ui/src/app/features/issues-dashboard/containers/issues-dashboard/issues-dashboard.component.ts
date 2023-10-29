import { Component, OnInit } from '@angular/core';
import { IssuesService } from 'app/features/issues-dashboard/services/issues.service';
import { IssueSearchCriteriaDTO } from 'app/shared/interfaces/IssueSearchCriteriaDTO';

@Component({
  selector: 'app-issues-dashboard',
  templateUrl: './issues-dashboard.component.html',
  styleUrls: [ './issues-dashboard.component.scss' ]
})
export class IssuesDashboardComponent implements OnInit {
  public containers: string[];

  constructor(private issuesService: IssuesService) {
  }

  public ngOnInit(): void {
    this.getIssuesContainers();
  }

  public searchIssues(searchCriteria: IssueSearchCriteriaDTO): void {

  }

  private getIssuesContainers(): void {
    this.issuesService.getIssuesContainers().subscribe((containers: string[]) => {
      this.containers = containers;
    });
  }
}
