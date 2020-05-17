import { Component, OnInit } from '@angular/core';
import { WeatherService } from '../weather.service';

@Component({
  selector: 'app-set-weather',
  templateUrl: './set-weather.component.html',
  styleUrls: ['./set-weather.component.css']
})
export class SetWeatherComponent implements OnInit {
  name: string;
  constructor(private weatherService: WeatherService) { }

  ngOnInit(): void {
  }

  onClick() {
    localStorage.setItem(this.weatherService.storageKey, this.name);
    this.weatherService.notifyCityUpdate();
  }
}
