import { Injectable } from '@angular/core';
import { HttpHandler, HttpInterceptor, HttpRequest } from '@angular/common/http';
import { ApplicationStateService } from 'app/shared/services/application-state.service';
import { catchError, finalize, Observable, throwError } from 'rxjs';

@Injectable()
export class LoadingInterceptor implements HttpInterceptor {

  constructor(private applicationStateService: ApplicationStateService) {
  }

  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<any> {
    this.applicationStateService.isLoading = true;
    return next.handle(req).pipe(
      finalize(() => this.applicationStateService.isLoading = false),
      catchError((err) => {
        this.applicationStateService.isLoading = false
        return throwError(err);
      })
    )
  }
}