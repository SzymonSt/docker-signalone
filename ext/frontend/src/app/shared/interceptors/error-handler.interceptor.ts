import { Injectable } from '@angular/core';
import { HttpHandler, HttpInterceptor, HttpRequest } from '@angular/common/http';
import { catchError, Observable, throwError } from 'rxjs';
import { ToastrService } from 'ngx-toastr';
import { TranslateService } from '@ngx-translate/core';
import { AuthStateService } from 'app/auth/services/auth-state.service';

@Injectable()
export class ErrorHandlerInterceptor implements HttpInterceptor {

  constructor(private toastrService: ToastrService, private translateService: TranslateService, private authStateService: AuthStateService) {
  }

  public intercept(req: HttpRequest<any>, next: HttpHandler): Observable<any> {
    if (req.url.includes('token/refresh')) {
      return next.handle(req).pipe(
        catchError((err) => {
          if (err.status === 401) {
            this.authStateService.logout();
            this.toastrService.error(this.translateService.instant('ERROR.UNAUTHORIZED'), this.translateService.instant('UI.ERROR'));

          } else {
            this.toastrService.error(this.translateService.instant(err.descriptionKey ? `ERROR.${err.descriptionKey}` : 'ERROR.UNDEFINED'), this.translateService.instant('UI.ERROR'));
          }
          return throwError(err);
        })
      );
    } else {
      return next.handle(req);
    }
  }
}