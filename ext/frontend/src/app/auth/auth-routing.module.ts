import { NgModule } from "@angular/core";
import { RouterModule, Routes } from "@angular/router";
import { LoginComponent } from 'app/auth/components/login/login.component';
import { NotLoggedInGuardService } from 'app/shared/guards/not-logged-in-guard.service';

const routes: Routes = [
  { path: "login", component: LoginComponent, canActivate: [NotLoggedInGuardService] },
];

@NgModule({
  imports: [ RouterModule.forChild(routes) ],
  exports: [ RouterModule ],
})
export class AuthRoutingModule {
}
