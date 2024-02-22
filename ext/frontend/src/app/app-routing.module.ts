import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { LoggedInGuardService } from 'app/shared/guards/logged-in-guard.service';

const routes: Routes = [
  {
    path: '',
    redirectTo: 'issues-dashboard',
    pathMatch: 'full'
  },
  {
    path: 'issues-dashboard',
    loadChildren: () => import('app/features/issues-dashboard/issues-dashboard.module').then(m => m.IssuesDashboardModule),
    canActivate: [LoggedInGuardService]
  }
];

@NgModule({
  imports: [ RouterModule.forRoot(routes) ],
  exports: [ RouterModule ]
})
export class AppRoutingModule {
}
