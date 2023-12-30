import { Injectable } from '@angular/core';
import { HttpHandler, HttpInterceptor, HttpRequest } from '@angular/common/http';
import { catchError, Observable, throwError } from 'rxjs';
import { ToastrService } from 'ngx-toastr';
import { TranslateService } from '@ngx-translate/core';

@Injectable()
export class ErrorHandlerInterceptor implements HttpInterceptor {

  constructor(private toastrService: ToastrService, private translateService: TranslateService,) {
  }

  public intercept(req: HttpRequest<any>, next: HttpHandler): Observable<any> {
    return next.handle(req).pipe(
      catchError((err) => {
        this.toastrService.error(this.translateService.instant(err.descriptionKey ? `ERROR.${err.descriptionKey}` : 'ERROR.UNDEFINED'), this.translateService.instant('UI.ERROR'));
        return throwError(err);
      })
    );
  }
}