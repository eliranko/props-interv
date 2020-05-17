import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Weather } from './models/Weather';
import { Observable, BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class WeatherService {
  storageKey: string = "cityName";
  private updatedCitySource = new BehaviorSubject('');
  public updatedCity = this.updatedCitySource.asObservable();

  constructor(private http: HttpClient) { }

  notifyCityUpdate() {
    this.updatedCitySource.next('');
  }

  getWeather(cityName: string): Observable<Weather> {
    return this.http.get<Weather>("api/weather/" + cityName);
  }
}
