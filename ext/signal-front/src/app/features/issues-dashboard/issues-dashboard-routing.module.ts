import { NgModule } from "@angular/core";
import { RouterModule, Routes } from "@angular/router";
import {
  IssuesDashboardComponent
} from 'app/features/issues-dashboard/containers/issues-dashboard/issues-dashboard.component';

const routes: Routes = [
  { path: "", component: IssuesDashboardComponent },
];

@NgModule({
  imports: [ RouterModule.forChild(routes) ],
  exports: [ RouterModule ],
})
export class IssuesDashboardRoutingModule {
}
