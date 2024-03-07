import { Component, EventEmitter, Input, Output } from '@angular/core';
import { DetailedIssueDTO, DetailedIssueScore } from 'app/shared/interfaces/DetailedIssueDTO';

@Component({
  selector: 'app-solutions-list',
  templateUrl: './solutions-list.component.html',
  styleUrls: [ './solutions-list.component.scss' ]
})
export class SolutionsListComponent {
  @Input()
  public activeIssue: DetailedIssueDTO;
  @Output()
  public scoreSelected: EventEmitter<DetailedIssueScore> = new EventEmitter<DetailedIssueScore>();
  @Output()
  public markIssueAsResolved: EventEmitter<void> = new EventEmitter<void>();
  @Output()
  public regenerateIssue: EventEmitter<void> = new EventEmitter<void>();
  constructor() {
  }

  public goToSolutionSource(url : string): void {
    window.open(url, "_blank")
  }

  public positiveScoreSelected(): void {
    if (this.activeIssue.score === 1) {
      this.activeIssue.score = 0;
    } else {
      this.activeIssue.score = 1;
    }
    this.scoreSelected.emit(this.activeIssue.score);
  }

  public negativeScoreSelected(): void {
    if (this.activeIssue.score === -1) {
      this.activeIssue.score = 0;
    } else {
      this.activeIssue.score = -1;
    }
    this.scoreSelected.emit(this.activeIssue.score);
  }
}
