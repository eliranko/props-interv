package main

type weather struct {
	Coord       corditations  `json:"coord"`
	WeatherData []weatherData `json:"weather"`
	MainData    mainData      `json:"main"`
	Name        string        `json:"name"`
	Response    int           `json:"cod"`
}

type corditations struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type weatherData struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type mainData struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	MinTemp   float64 `json:"temp_min"`
	MaxTemp   float64 `json:"temp_max"`
	Humidity  float64 `json:"humidity"`
}
