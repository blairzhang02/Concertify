import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ConcertService {

  private apiUrl = 'http://localhost:8080';  // URL to your Go web server

  constructor(private http: HttpClient) { }

  // Method to start the authentication process
  initiateAuth(): Observable<any> {
    // This should match the endpoint that starts the OAuth process in your Go app
    return this.http.get(`${this.apiUrl}/auth`);
  }

  // Method to retrieve concerts (or any other data your Go backend provides)
  getConcerts(): Observable<any> {
    // The URL here would depend on your Go application's endpoints
    return this.http.get(`${this.apiUrl}/concerts`);
  }
}
