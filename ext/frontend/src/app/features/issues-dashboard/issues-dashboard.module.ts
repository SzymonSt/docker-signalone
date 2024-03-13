import { NgModule } from '@angular/core';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { CommonModule } from '@angular/common';
import { SvgIconComponent } from 'angular-svg-icon';
import { SharedModule } from 'app/shared/SharedModule';
import {
  IssuesDashboardComponent
} from 'app/features/issues-dashboard/containers/issues-dashboard/issues-dashboard.component';
import { IssuesDashboardRoutingModule } from 'app/features/issues-dashboard/issues-dashboard-routing.module';
import {
  IssuesSearchCriteriaComponent
} from 'app/features/issues-dashboard/components/left-panel-components/issues-search-criteria/issues-search-criteria.component';
import {
  IssuesRightPanelComponent
} from 'app/features/issues-dashboard/components/issues-right-panel/issues-right-panel.component';
import {
  IssuesLeftPanelComponent
} from 'app/features/issues-dashboard/components/issues-left-panel/issues-left-panel.component';
import {
  IssuesListComponent
} from 'app/features/issues-dashboard/components/left-panel-components/issues-list/issues-list.component';
import { NgSelectModule } from '@ng-select/ng-select';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { BsDatepickerModule } from 'ngx-bootstrap/datepicker';
import {
  IssueCellComponent
} from 'app/features/issues-dashboard/components/left-panel-components/issue-cell/issue-cell.component';
import { PaginationModule } from 'ngx-bootstrap/pagination';
import {
  SolutionsChatComponent
} from 'app/features/issues-dashboard/components/right-panel-components/chat/solutions-chat.component';
import {
  SolutionsListComponent
} from 'app/features/issues-dashboard/components/right-panel-components/solutions-list/solutions-list.component';
import { MarkdownModule } from 'ngx-markdown';

@NgModule({
  declarations: [ IssuesDashboardComponent, IssuesSearchCriteriaComponent, IssuesRightPanelComponent, IssuesLeftPanelComponent, IssuesListComponent, IssueCellComponent, SolutionsChatComponent, SolutionsListComponent ],
  imports: [
    CommonModule,
    TranslateModule,
    SharedModule,
    IssuesDashboardRoutingModule,
    NgSelectModule,
    FormsModule,
    ReactiveFormsModule,
    NgbModule,
    BsDatepickerModule,
    PaginationModule,
    MarkdownModule.forRoot(),
    MatTooltipModule,
    SvgIconComponent
  ],
  exports: []
})
export class IssuesDashboardModule {
}
