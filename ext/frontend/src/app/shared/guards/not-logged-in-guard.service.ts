import { Injectable } from '@angular/core';
import { Router, CanActivate } from '@angular/router';
import { AuthStateService } from 'app/auth/services/auth-state.service';

@Injectable({ providedIn: 'root' })
export class NotLoggedInGuardService implements CanActivate {
  constructor(public authStateService: AuthStateService, public router: Router) {}
  canActivate(): boolean {
    console.log(this.authStateService.isLoggedIn)
  if (this.authStateService.isLoggedIn) {
    this.router.navigate(['issues-dashboard']);
    return false;
  }
  return true;
}}