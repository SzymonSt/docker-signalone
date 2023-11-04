import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

const routes: Routes = [
  {
    path: '',
    redirectTo: 'issues-dashboard',
    pathMatch: 'full'
  },
  {
    path: 'issues-dashboard',
    loadChildren: () => import('app/features/issues-dashboard/issues-dashboard.module').then(m => m.IssuesDashboardModule)
  }
];

@NgModule({
  imports: [ RouterModule.forRoot(routes) ],
  exports: [ RouterModule ]
})
export class AppRoutingModule {
}
