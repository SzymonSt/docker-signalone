import { Component, OnInit } from '@angular/core';
import { IssuesService } from 'app/features/issues-dashboard/services/issues.service';
import { IssueSearchCriteriaDTO } from 'app/shared/interfaces/IssueSearchCriteriaDTO';
import { IssueDTO } from 'app/shared/interfaces/IssueDTO';
import { DetailedIssueDTO } from 'app/shared/interfaces/DetailedIssueDTO';

@Component({
  selector: 'app-issues-dashboard',
  templateUrl: './issues-dashboard.component.html',
  styleUrls: [ './issues-dashboard.component.scss' ]
})
export class IssuesDashboardComponent implements OnInit {
  public containers: string[];
  public issues: IssueDTO[];
  public activeIssue: DetailedIssueDTO;
  public max: number;
  public activePage: number = 1;
  public isSidebarHidden: boolean = false;
  private lastSearchCriteria: IssueSearchCriteriaDTO = new IssueSearchCriteriaDTO();
  constructor(private issuesService: IssuesService) {
  }

  public ngOnInit(): void {
    this.getIssuesContainers();
    this.searchIssues(this.lastSearchCriteria);
    this.subscribeIssuesContainers();
  }

  public searchIssues(searchCriteria?: IssueSearchCriteriaDTO, revokeLoader: boolean = false): void {
    if (searchCriteria) {
      this.activePage = searchCriteria.offset ? searchCriteria.offset * searchCriteria.limit : 1;
      this.lastSearchCriteria = {
        ...this.lastSearchCriteria,
        ...searchCriteria
      };
    }
    this.issuesService.getIssuesList(this.lastSearchCriteria, revokeLoader).subscribe((res) => {
      this.issues = res.issues
      this.max = res.max;
    });
  }

  public viewIssue(issue: IssueDTO): void {
    this.issuesService.getIssue(issue.id).subscribe((response) => {
      this.activeIssue = response;
    });

  }

  private getIssuesContainers(): void {
    this.issuesService.getIssuesContainers().subscribe((containers: string[]) => {
      this.containers = containers;
    });
  }

  private subscribeIssuesContainers(): void {
    setInterval(() => {
      this.searchIssues(this.lastSearchCriteria, true);
    }, 15000)
  }
}
