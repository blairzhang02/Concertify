import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AuthService } from '../auth.service'; // if you have a service for handling authentication

@Component({
  selector: 'app-callback',
  templateUrl: './callback.component.html',
  styleUrls: ['./callback.component.css']
})
export class CallbackComponent implements OnInit {

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private authService: AuthService // Injecting our own AuthService
  ) { }

  ngOnInit() {
    this.route.fragment.subscribe((fragment: string | null) => { // here you accept that fragment could be string or null
      if (typeof fragment === 'string') { // checking if fragment is a string
        const params = new URLSearchParams(fragment);
        const accessToken = params.get('access_token');
        const error = params.get('error');
  
        if (accessToken) {
          // Save your access token in AuthService or local storage
          this.authService.setAuthToken(accessToken);
          
          // Redirect to another page, e.g., your main application page
          this.router.navigate(['/']);
        } else if (error) {
          // Handle errors
          console.error(error);
          // Redirect to an error page or show an error message
        }
      } else {
        // Potentially handle the case where fragment is null
        console.error('No fragment found in the callback URL.');
      }
    });
  }
  
}
