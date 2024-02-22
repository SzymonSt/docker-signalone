import { Injectable } from '@angular/core';
import { Router, CanActivate } from '@angular/router';
import { AuthStateService } from 'app/auth/services/auth-state.service';

@Injectable({ providedIn: 'root' })
export class LoggedInGuardService implements CanActivate {
  constructor(public authStateService: AuthStateService, public router: Router) {}
  canActivate(): boolean {
    if (!this.authStateService.isLoggedIn) {
      this.router.navigate(['login']);
      return false;
    }
    return true;
}}