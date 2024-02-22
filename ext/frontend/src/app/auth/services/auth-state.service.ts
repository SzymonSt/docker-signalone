import { Injectable, NgZone, OnDestroy } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import * as _ from 'lodash';
import * as moment from 'moment';
import { Duration } from 'moment';
import { Token } from 'app/shared/interfaces/Token';
import { AuthService } from 'app/auth/services/auth.service';
import { Router } from '@angular/router';
import { SocialUser } from '@abacritt/angularx-social-login';

@Injectable({ providedIn: 'root' })
export class AuthStateService implements OnDestroy {
  private static readonly TOKEN_REFRESH_INTERVAL: Duration = moment.duration('1', 'minutes');
  public token$: BehaviorSubject<Token> = new BehaviorSubject<Token>(null);
  public isLoggedIn$: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  private wasNavigatedFromLogin: boolean = false;
  private tokenRefreshIntervalId!: ReturnType<typeof setInterval>;
  constructor(private zone: NgZone,
              private authService: AuthService,
              private router: Router,) {}

  public get token(): Token {
    return this.token$.value;
  }

  public set token(value: Token) {
    this.token$.next(value);
  }

  public get isLoggedIn(): boolean {
    return this.isLoggedIn$.value;
  }

  public set isLoggedIn(value: boolean) {
    this.isLoggedIn$.next(value);
  }

  public ngOnDestroy(): void {
    this.cancelTokenRefreshSchedule();
  }

  public login(email: string, password: string, silent: boolean = false): Promise<Token> {
    return new Promise((resolve, reject) => {
      this.authService.login(email, password).toPromise()
        .then((result: { token: Token }) => {
          this.setToken(result.token)
            .then((savedToken: Token) => {
              this.manageLoginSuccess(result)
              resolve(this.token);
            })
            .catch((error) => {
              this.token = null;
              this.isLoggedIn = false;
              reject(error);
            });
        })
        .catch((error: any) => {
          this.token = null;
          this.isLoggedIn = false;
          reject(error);
        });
    });
  }

  public loginWithGoogle(user: SocialUser): Promise<Token> {
    return new Promise((resolve, reject) => {
      this.authService.loginWithGoogle(user).toPromise()
        .then((result: { token: Token }) => {
          this.setToken(result.token)
            .then((savedToken: Token) => {
              this.manageLoginSuccess(result)
              resolve(this.token);
            })
            .catch((error) => {
              this.token = null;
              this.isLoggedIn = false;
              reject(error);
            });
        })
        .catch((error: any) => {
          this.token = null;
          this.isLoggedIn = false;
          reject(error);
        });
    });
  }

  public loginWithGithub(code: string): Promise<Token> {
    return new Promise((resolve, reject) => {
      this.authService.loginWithGithub(code).toPromise()
        .then((result: { token: Token }) => {
          this.setToken(result.token)
            .then((savedToken: Token) => {
              this.manageLoginSuccess(result)
              resolve(this.token);
            })
            .catch((error) => {
              this.token = null;
              this.isLoggedIn = false;
              reject(error);
            });
        })
        .catch((error: any) => {
          this.token = null;
          this.isLoggedIn = false;
          reject(error);
        });
    });
  }

  public logout(silent: boolean = false): Promise<void> {
    return new Promise((resolve, reject) => {
      this.authService.logout(this.token).toPromise()
        .then(() => {
          this.deleteToken()
            .then(() => {
              this.token = null;
              this.isLoggedIn = false;
              this.cancelTokenRefreshSchedule();
              this.goToLogin();
              resolve();
            })
            .catch((error) => {
              // shouldn't happen really, but we're covering for safety
              // we need to wipe the state anyway, even if there's an error
              this.token = null;
              this.isLoggedIn = false;
              this.cancelTokenRefreshSchedule();
              reject(error);
            });
        })
        .catch((error: any) => {
          // we need to wipe the state anyway, even if there's an error
          // we don't care about the result of the deleteToken operation at this point, we're in error state anyway
          this.deleteToken()
            .finally(() => {
              this.token = null;
              this.isLoggedIn = false;
              this.cancelTokenRefreshSchedule();
              this.goToLogin();
              reject(error);
            });
        });
    });
  }

  public refresh(token: Token): Promise<Token> {
    return new Promise((resolve, reject) => {
      this.authService.refresh(token).toPromise()
        .then((result: { token: Token }) => {
          this.setToken(result.token)
            .then((savedToken: Token) => {
              this.token = result.token;
              this.rescheduleRefresh(this.token);
              resolve(this.token);
            })
            .catch((error) => {
              reject(error);
            });
        })
        .catch((error: any) => {
          reject(error);
        });
    });
  }

  public recoverToken(): Promise<Token> {
    return new Promise<Token>((resolve, reject) => {
      this.getToken()
        .then((token: Token) => {
          if (_.isNil(token)) {
            this.isLoggedIn = false;
            reject();
          } else {
            if (token.isExpired() || token.isNearlyExpired()) {
              this.refresh(token)
                .then((refreshedToken: Token) => {
                  this.isLoggedIn = true;

                  resolve(refreshedToken);
                })
                .catch((error) => {
                  this.isLoggedIn = false;
                  reject(error);
                });
            } else {
              this.token = token;
              this.isLoggedIn = true;

              this.scheduleTokenRefresh(token);

              resolve(token);
            }
          }
        })
        .catch((error) => {
          this.isLoggedIn = false;
          reject(error);
        });
    });
  }

  // TODO Add to be called after unauthorized
  public clearTokenData(): void {
    this.deleteToken()
      .finally(() => {
        this.token = null;
        this.cancelTokenRefreshSchedule();
      });
  }

  private scheduleTokenRefresh(token: Token): void {
    this.zone.runOutsideAngular(() => {
      this.tokenRefreshIntervalId = setInterval(() => {
        this.zone.run(() => {
          if (this.token.isNearlyExpired()) {
            this.refresh(this.token)
              .then((refreshedToken: Token) => {
              })
              .catch((error) => {
              });
          }
        });
      }, AuthStateService.TOKEN_REFRESH_INTERVAL.as('milliseconds'));
    });
  }

  private cancelTokenRefreshSchedule(): void {
    if (this.tokenRefreshIntervalId) {
      clearInterval(this.tokenRefreshIntervalId);
      // @ts-ignore
      this.tokenRefreshIntervalId = null;
    }
  }
  private rescheduleRefresh(token: Token): void {
    this.cancelTokenRefreshSchedule();
    this.scheduleTokenRefresh(token);
  }

  private setToken(token: Token): Promise<Token> {
    return this.authService.setToken(token).toPromise();
  }

  private getToken(): Promise<Token> {
    return this.authService.getToken().toPromise();
  }

  private deleteToken(): Promise<void> {
    return this.authService.deleteToken().toPromise();
  }

  private goToDashboard():void {
    if (!this.wasNavigatedFromLogin) {
      this.router.navigateByUrl('/issues-dashboard')
    }
    this.wasNavigatedFromLogin = true;
  }
  private goToLogin():void {
    this.router.navigateByUrl('/login')
    this.wasNavigatedFromLogin = false;
  }

  private manageLoginSuccess(result: {token: Token}): void {
    this.token = result.token;
    this.isLoggedIn = true;
    if (!_.isNil(this.token)) {
      this.scheduleTokenRefresh(this.token);
    }

    this.goToDashboard();
  }
}
