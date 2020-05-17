import { Component, OnInit } from '@angular/core';
import { MoviesService } from '../movies.service';
import { Movie } from '../models/Movie';

@Component({
  selector: 'app-movie-details',
  templateUrl: './movie-details.component.html',
  styleUrls: ['./movie-details.component.css']
})
export class MovieDetailsComponent implements OnInit {
  name: string;
  chosenMovie: Movie;

  constructor(private movies: MoviesService) { }

  ngOnInit(): void {
  }

  onSearch() {
    this.movies.getMovie(this.name).subscribe(movie => {
      this.chosenMovie = movie
      console.log(this.chosenMovie);
    });
  }
}
