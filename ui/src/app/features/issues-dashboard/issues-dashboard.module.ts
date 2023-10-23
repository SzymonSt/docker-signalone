import { NgModule } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { CommonModule } from '@angular/common';
import { SharedModule } from 'app/shared/SharedModule';
import {
  IssuesDashboardComponent
} from 'app/features/issues-dashboard/containers/issues-dashboard/issues-dashboard.component';
import { IssuesDashboardRoutingModule } from 'app/features/issues-dashboard/issues-dashboard-routing.module';
import {
  IssuesSearchCriteriaComponent
} from 'app/features/issues-dashboard/components/issues-search-criteria/issues-search-criteria.component';
import {
  IssuesRightPanelComponent
} from 'app/features/issues-dashboard/components/issues-right-panel/issues-right-panel.component';
import {
  IssuesLeftPanelComponent
} from 'app/features/issues-dashboard/components/issues-left-panel/issues-left-panel.component';
import { IssuesListComponent } from 'app/features/issues-dashboard/components/issues-list/issues-list.component';
import { NgSelectModule } from '@ng-select/ng-select';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { BsDatepickerModule } from 'ngx-bootstrap/datepicker';

@NgModule({
  declarations: [ IssuesDashboardComponent, IssuesSearchCriteriaComponent, IssuesRightPanelComponent, IssuesLeftPanelComponent, IssuesListComponent ],
  imports: [
    CommonModule,
    TranslateModule,
    SharedModule,
    IssuesDashboardRoutingModule,
    NgSelectModule,
    FormsModule,
    ReactiveFormsModule,
    NgbModule,
    BsDatepickerModule
  ],
  exports: []
})
export class IssuesDashboardModule {
}
