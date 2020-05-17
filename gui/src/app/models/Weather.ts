export class Weather {
    coord: Coordinations;
    weather: WeatherData[];
    main: Main;
    name: string;
}

export class Coordinations {
    lon: number;
    lat: number;
}

export class WeatherData {
    id: number;
    main: string;
    description: string;
    icon: string;
}

export class Main {
    temp: number;
    feels_like: number;
    temp_min: number;
    temp_max: number;
    humidity: number;
}