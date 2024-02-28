import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { OAuth2TokenDTO } from 'app/shared/interfaces/OAuth2TokenDTO';
import { Token } from 'app/shared/interfaces/Token';
import { StorageUtil } from 'app/shared/util/StorageUtil';
import { environment } from 'environment/environment.development';
import { map, Observable } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private static readonly TOKEN_KEY: string = 'signal_token';
  constructor(private httpClient: HttpClient, private storageUtil: StorageUtil) {
  }

  public login(email: string, password: string): Observable<{ token: Token }> {
    return this.httpClient.post<{ token: Token }>(`${environment.authUrl}/login`, { email, password })
      .pipe(
        map((response: any) => {
          const token: OAuth2TokenDTO = OAuth2TokenDTO.fromOAuth2Object(response);
          return { token: token};
        })
      );
  }

  public loginWithGoogle(accessToken: string): Observable<{ token: Token }> {
    return this.httpClient.post<{ token: Token }>(`${environment.authUrl}/login-with-google`, { idToken: accessToken })
      .pipe(
        map((response: any) => {
          const token: OAuth2TokenDTO = OAuth2TokenDTO.fromOAuth2Object(response);
          return { token: token};
        })
      );
  }

  public loginWithGithub(code: string): Observable<{ token: Token }> {

    return this.httpClient.post<{ token: Token }>(`${environment.authUrl}/login-with-github`, {code})
      .pipe(
        map((response: any) => {
          const token: OAuth2TokenDTO = OAuth2TokenDTO.fromOAuth2Object(response);
          return { token: token};
        })
      );
  }

  public logout(token: Token): Observable<void> {
    return this.httpClient.post(`${environment.authUrl}/logout`, { refreshToken: token.refreshToken })
      .pipe(
        map(() => {
          return;
        })
      );
  }

  public refresh(token: Token): Observable<{ token: Token }> {
    const headers: HttpHeaders = new HttpHeaders({
      'Content-Type': 'application/json',
      Accept: 'application/json'
    });

    return this.httpClient.post(`${environment.authUrl}/token/refresh`, JSON.stringify({ refreshToken: token.refreshToken }), { headers })
      .pipe(
        map((response: any) => {
          const refreshedToken: OAuth2TokenDTO = OAuth2TokenDTO.fromOAuth2Object(response);
          return { token: refreshedToken};
        })
      );
  }

  public startPasswordReset(email: string): Observable<void> {
    return this.httpClient.post(`${environment.apiUrl}/accounts/${encodeURIComponent(email)}/password/init-reset`, {})
      .pipe(
        map(() => {
          return;
        })
      );
  }

  public setToken(token: Token): Observable<OAuth2TokenDTO> {
    return new Observable<OAuth2TokenDTO>((observer) => {
      this.storageUtil.saveData<Token>(token, AuthService.TOKEN_KEY)
        .then((result: any) => {
          observer.next(result);
          observer.complete();
        })
        .catch((error: any) => {
          observer.error(error);
        });
    });
  }

  public getToken(): Observable<Token> {
    return new Observable<Token>((observer) => {
      this.storageUtil.loadData<OAuth2TokenDTO>(AuthService.TOKEN_KEY, OAuth2TokenDTO)
        .then((result: any) => {
          observer.next(result);
          observer.complete();
        })
        .catch((error: any) => {
          observer.error(error);
        });
    });
  }

  public deleteToken(): Observable<void> {
    return new Observable<void>((observer) => {
      this.storageUtil.deleteData(AuthService.TOKEN_KEY)
        .then(() => {
          observer.next();
          observer.complete();
        })
        .catch((error: any) => {
          observer.error(error);
        });
    });
  }

}