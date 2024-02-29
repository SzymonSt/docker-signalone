import { Component, OnInit } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { AuthStateService } from 'app/auth/services/auth-state.service';
import { IssuesService } from 'app/features/issues-dashboard/services/issues.service';
import { DetailedIssueDTO, DetailedIssueScore } from 'app/shared/interfaces/DetailedIssueDTO';
import { IssueDTO } from 'app/shared/interfaces/IssueDTO';
import { IssueSearchCriteriaDTO } from 'app/shared/interfaces/IssueSearchCriteriaDTO';
import { RateIssueDTO } from 'app/shared/interfaces/RateIssueDTO';

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
  private issuesContainersRefreshInterval: any;
  constructor(private issuesService: IssuesService, private authStateService: AuthStateService) {
    this.getIssuesContainers();
    this.authStateService.isLoggedIn$.pipe(takeUntilDestroyed()).subscribe(isLoggedIn => {
      if (isLoggedIn) {
        this.subscribeIssuesContainers();
      } else {
        this.clearIssuesContainersSubscription();
      }
    })
  }

  public ngOnInit(): void {
    this.searchIssues(this.lastSearchCriteria);
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

  public scoreSelected(score: DetailedIssueScore): void {
    this.issuesService.rateIssue(this.activeIssue.id, new RateIssueDTO(score)).subscribe()
  }

  private getIssuesContainers(): void {
    this.issuesService.getIssuesContainers().subscribe((containers: string[]) => {
      this.containers = containers;
    });
  }

  private subscribeIssuesContainers(): void {
    this.clearIssuesContainersSubscription();
    this.issuesContainersRefreshInterval = setInterval(() => {
      this.searchIssues(this.lastSearchCriteria, true);
    }, 15000)
  }

  private clearIssuesContainersSubscription():void {
    if (this.issuesContainersRefreshInterval) {
      clearInterval(this.issuesContainersRefreshInterval)
      this.issuesContainersRefreshInterval = null
    }
  }
}
