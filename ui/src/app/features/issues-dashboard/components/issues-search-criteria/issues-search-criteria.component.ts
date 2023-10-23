import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { IssueType } from 'app/shared/enum/IssueType';
import { Severity } from 'app/shared/enum/Severity';
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
  public severityOptions: Severity[] = Object.values(Severity);
  public todayDate: Date = new Date();
  public searchForm: FormGroup;
  public isSubmitted = false;

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
  }

  private initForm(): void {
    this.searchForm = new FormGroup({
      searchText: new FormControl(null),
      type: new FormControl(null),
      severity: new FormControl(null),
      container: new FormControl(null),
      dateFrom: new FormControl(null),
      dateTo: new FormControl(null),
      showOnlyUnresolved: new FormControl(null),
    }, { validators: dateRangeValidator() })
  }
}
