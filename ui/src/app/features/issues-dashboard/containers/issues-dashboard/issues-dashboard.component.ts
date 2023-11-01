import { Component, OnInit } from '@angular/core';
import { IssuesService } from 'app/features/issues-dashboard/services/issues.service';
import { IssueSearchCriteriaDTO } from 'app/shared/interfaces/IssueSearchCriteriaDTO';
import { IssueDTO } from 'app/shared/interfaces/IssueDTO';

@Component({
  selector: 'app-issues-dashboard',
  templateUrl: './issues-dashboard.component.html',
  styleUrls: [ './issues-dashboard.component.scss' ]
})
export class IssuesDashboardComponent implements OnInit {
  public containers: string[];
  public issues: IssueDTO[];
  public activeIssue: IssueDTO;

  constructor(private issuesService: IssuesService) {
  }

  public ngOnInit(): void {
    this.getIssuesContainers();
    this.searchIssues();
  }

  public searchIssues(searchCriteria?: IssueSearchCriteriaDTO): void {
    this.issuesService.getIssuesList(searchCriteria).subscribe((issues) => {
      this.issues = issues;
    });
  }

  public viewIssue(issue: IssueDTO): void {
    console.log('ACTIVE ISSUE', issue)
    this.activeIssue = issue;
  }

  private getIssuesContainers(): void {
    this.issuesService.getIssuesContainers().subscribe((containers: string[]) => {
      this.containers = containers;
    });
  }
}
