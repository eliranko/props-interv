import { Component, OnInit, Input } from '@angular/core';
import { Weather } from '../models/Weather';

@Component({
  selector: 'app-inline-weather-data',
  templateUrl: './inline-weather-data.component.html',
  styleUrls: ['./inline-weather-data.component.css']
})
export class InlineWeatherDataComponent implements OnInit {
  @Input() weather: Weather;
  constructor() { }

  ngOnInit(): void {
  }

}
