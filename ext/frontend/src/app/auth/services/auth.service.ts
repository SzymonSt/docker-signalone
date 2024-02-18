import { Injectable } from '@angular/core';
import { Token } from 'app/shared/interfaces/Token';
import { map, Observable } from 'rxjs';
import { OAuth2TokenDTO } from 'app/shared/interfaces/OAuth2TokenDTO';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { environment } from 'environment/environment.development';
import { StorageUtil } from 'app/shared/util/StorageUtil';
import { SocialUser } from '@abacritt/angularx-social-login';
import { HttpEncoder } from 'app/shared/util/HttpEncoder';
import { NormalizeObjectValue } from 'app/shared/util/NormalizeObjectValue';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private static readonly TOKEN_KEY: string = 'token';
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

  public loginWithGoogle(user: SocialUser): Observable<{ token: Token }> {
    return this.httpClient.post<{ token: Token }>(`${environment.authUrl}/login-with-google`, { user})
      .pipe(
        map((response: any) => {
          const token: OAuth2TokenDTO = OAuth2TokenDTO.fromOAuth2Object(response);
          return { token: token};
        })
      );
  }

  public loginWithGithub(): Observable<{ token: Token }> {
    const params: HttpParams = new HttpParams().set('client_id', 'c88a4f9d4868879974e').set('redirect_uri', 'http://localhost:8080/oauth/redirect');
    return this.httpClient.get<{ token: Token }>(`https://github.com/login/oauth/authorize`, {params})
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

  public completePasswordReset(email: string, verificationCode: string, newPassword: string): Observable<void> {
    const request: {
      verificationCode: string;
      newPassword: string;
    } = {
      verificationCode: verificationCode,
      newPassword: newPassword
    };

    return this.httpClient.post(`${environment.apiUrl}/accounts/${encodeURIComponent(email)}/password/set-new`, request)
      .pipe(
        map(() => {
          return;
        })
      );
  }

  public changePassword(currentPassword: string, newPassword: string): Observable<void> {
    const request: {
      oldPassword: string;
      newPassword: string;
    } = {
      oldPassword: currentPassword,
      newPassword: newPassword
    };

    return this.httpClient.patch<void>(`${environment.apiUrl}/accounts/me/password`, request);
  }

  public changePasswordForced(newPassword: string): Observable<void> {
    const request: {
      newPassword: string;
    } = {
      newPassword: newPassword
    };

    return this.httpClient.put<void>(`${environment.apiUrl}/users/me/passwordForced`, request);
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