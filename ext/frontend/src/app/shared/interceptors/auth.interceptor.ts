import { Injectable } from '@angular/core';
import { HttpEvent, HttpHandler, HttpInterceptor, HttpRequest } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from 'environment/environment';
import { AuthStateService } from 'app/auth/services/auth-state.service';

@Injectable({ providedIn: 'root' })
export class AuthInterceptor implements HttpInterceptor {

  constructor(private authStateService: AuthStateService) {
  }

  public intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    if (this.authStateService.token && request.method && (
      (request.url?.indexOf(environment.apiUrl) > -1)
    )) {
      const authorizedRequest: HttpRequest<any> = request.clone({
        setHeaders: {
          Authorization: 'Bearer ' + this.authStateService.token.accessToken
        }
      });

      return next.handle(authorizedRequest);
    } else {
      return next.handle(request);
    }
  }
}
