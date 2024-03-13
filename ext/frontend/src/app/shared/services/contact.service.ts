import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { ContactRequestDTO } from 'app/shared/interfaces/ContactRequestDTO';
import { environment } from 'environment/environment';

@Injectable({ providedIn: 'root' })
export class ContactService {

  constructor(private httpClient: HttpClient) {
  }

  public sendContactMessage(contactRequestBody: ContactRequestDTO) {
    return this.httpClient.post<void>(`${environment.apiUrl}/contact`, contactRequestBody)
  }
}