import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { IssueType } from 'app/shared/enum/IssueType';
import { IssueSeverity } from 'app/shared/enum/IssueSeverity';
import { FormControl, FormGroup } from '@angular/forms';
import { IssueSearchCriteriaDTO } from 'app/shared/interfaces/IssueSearchCriteriaDTO';
import { dateRangeValidator } from 'app/shared/validators/date-range.validator';

@Component({
  selector: 'app-issues-search-criteria',
  templateUrl: './issues-search-criteria.component.html',
  styleUrls: [ './issues-search-criteria.component.scss' ]
})
export class IssuesSearchCriteriaComponent implements OnInit {
  public issueTypeOptions: IssueType[] = Object.values(IssueType);
  public severityOptions: IssueSeverity[] = Object.values(IssueSeverity);
  public todayDate: Date = new Date();
  public searchForm: FormGroup;
  public isSubmitted = false;

  @Input()
  public containers: string[];

  @Output()
  public criteriaChanged: EventEmitter<IssueSearchCriteriaDTO> = new EventEmitter<IssueSearchCriteriaDTO>();

  constructor() {
  }

  public ngOnInit(): void {
    this.initForm();
  }

  public submitForm(): void {
    this.isSubmitted = true;
    this.searchForm.markAsDirty();
    this.searchForm.markAllAsTouched();
    if (this.searchForm.valid) {
      this.criteriaChanged.emit(this.searchForm.value);
    }
  }

  public clearForm(): void {
    this.searchForm.reset();
    this.criteriaChanged.emit(this.searchForm.value);
  }

  private initForm(): void {
    this.searchForm = new FormGroup({
      searchString: new FormControl(null),
      issueType: new FormControl(null),
      issueSeverity: new FormControl(null),
      containerId: new FormControl(null),
      startTimestamp: new FormControl(null),
      endTimestamp: new FormControl(null),
      isResolved: new FormControl(null),
    }, { validators: dateRangeValidator('startTimestamp', 'endTimestamp') })
  }
}
