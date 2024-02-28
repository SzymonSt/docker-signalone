import { Component, OnInit } from '@angular/core';
import { AuthStateService } from 'app/auth/services/auth-state.service';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-github-login',
  templateUrl: './google-login.component.html',
  styleUrls: [ './google-login.component.scss' ]
})
export class GoogleLoginComponent implements OnInit{

  constructor(private authStateService: AuthStateService, private activatedRoute: ActivatedRoute) {

  }

  public ngOnInit(): void {
    const idTokenFullString = this.activatedRoute.snapshot.fragment?.split('&')?.find(str => str?.includes('id_token'))?.split('=')[1]
    this.authStateService.loginWithGoogle(idTokenFullString);
  }

}
