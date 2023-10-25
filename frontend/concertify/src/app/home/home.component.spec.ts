// src/app/home/home.component.ts

import { Component, OnInit } from '@angular/core';
import { ConcertService } from '../concert.service'; // Import the service

// Make sure you've defined the Concert interface (you might have it from earlier steps)
export interface Concert {
  Artist: string;
  Venue: string;
  Date: string; // this could be 'Date' type if you handle date objects
  // ... other properties ...
}

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {
  concerts: Concert[] = []; // Holds the concert data

  constructor(private concertService: ConcertService) { }

  ngOnInit(): void {
    this.fetchConcerts();
  }

  fetchConcerts(): void {
    // Fetch the concerts from the service
    this.concertService.getConcerts().subscribe(
      (data: Concert[]) => {
        this.concerts = data; // Assign the data to our component's property
      },
      error => {
        console.error('There was an error fetching the concerts!', error);
      }
    );
  }
}
