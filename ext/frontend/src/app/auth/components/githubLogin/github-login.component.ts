import { Component, OnInit } from '@angular/core';
import { AuthStateService } from 'app/auth/services/auth-state.service';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-github-login',
  templateUrl: './github-login.component.html',
  styleUrls: [ './github-login.component.scss' ]
})
export class GithubLoginComponent implements OnInit{

  constructor(private authStateService: AuthStateService, private activatedRoute: ActivatedRoute) {

  }

  public ngOnInit(): void {
    this.authStateService.loginWithGithub(this.activatedRoute.snapshot.queryParams['code']);
  }

}
