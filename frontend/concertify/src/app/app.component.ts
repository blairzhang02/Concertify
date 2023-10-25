import { Component } from '@angular/core';
import { ConcertService } from './concert.service';
// This could be in a new file, e.g., `concert.model.ts`
export interface Concert {
  Artist: string;
  Venue: string;
  Date: string; // Or use 'Date' type if the date is an actual Date object
  // ... any other properties that a concert object might have
}

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  title = 'my-concert-finder';
  concerts: Concert[] = []; // Here you're specifying that 'concerts' will be an array of 'Concert' objects

  constructor(private concertService: ConcertService) {}

  ngOnInit() {
    this.startAuthProcess();
  }

  startAuthProcess(): void {
    // It's a good idea to retrieve this URL from your environment configuration or service.
    window.location.href = 'http://localhost:8080/auth'; // This URL should point to your backend auth initiation route.
  }
  

  fetchConcerts(): void {
    this.concertService.getConcerts().subscribe(concerts => {
      this.concerts = concerts; // TypeScript now knows what properties to expect on the objects in 'concerts'
    });
  }
}
