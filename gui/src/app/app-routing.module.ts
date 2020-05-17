import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { MovieDetailsComponent } from './movie-details/movie-details.component';
import { SetWeatherComponent } from './set-weather/set-weather.component';


const routes: Routes = [
  { path: 'movie', component: MovieDetailsComponent },
  { path: 'weather', component: SetWeatherComponent },
  { path: '', redirectTo: '/movie', pathMatch: 'full' },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
