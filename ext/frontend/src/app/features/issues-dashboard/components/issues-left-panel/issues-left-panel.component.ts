import { Component, EventEmitter, Input, Output } from '@angular/core';
import { IssueSearchCriteriaDTO } from 'app/shared/interfaces/IssueSearchCriteriaDTO';
import { IssueDTO } from 'app/shared/interfaces/IssueDTO';
import { PageChangedEvent } from 'ngx-bootstrap/pagination';
import { Constants } from 'app/config/Constant';

@Component({
  selector: 'app-issues-left-panel',
  templateUrl: './issues-left-panel.component.html',
  styleUrls: [ './issues-left-panel.component.scss' ]
})
export class IssuesLeftPanelComponent {
  @Input()
  public containers: string[];
  @Input()
  public issues: IssueDTO[];
  @Input()
  public max: number;
  @Input()
  public internalPage: number = 1;
  @Input()
  public isSidebarHidden: boolean;

  public pageSize: number = Constants.paginationLimit;

  @Output()
  public criteriaChanged: EventEmitter<IssueSearchCriteriaDTO> = new EventEmitter<IssueSearchCriteriaDTO>();
  @Output()
  public viewIssue: EventEmitter<IssueDTO> = new EventEmitter<IssueDTO>()
  @Output()
  public toggleSidebarVisibility: EventEmitter<boolean> = new EventEmitter<boolean>()

  private criteria: IssueSearchCriteriaDTO = new IssueSearchCriteriaDTO();

  constructor() {
  }

  public onPageChanged(event: PageChangedEvent): void {
    const newPage: number = event.page;
    const newPageIndex: number = newPage - 1;

    this.criteria = {
      ...this.criteria,
      limit: this.pageSize,
      offset: newPageIndex ? this.pageSize * newPageIndex : 0
    }
    this.criteriaChanged.emit(this.criteria);
  }

  public onCriteriaChange(newCriteria: IssueSearchCriteriaDTO): void {
    this.criteria = {
      ...this.criteria,
      ...newCriteria
    }
    this.criteriaChanged.emit(this.criteria);
  }

  public toggleSidebarHidden(): void {
    this.toggleSidebarVisibility.emit(!this.isSidebarHidden);
  }
}
