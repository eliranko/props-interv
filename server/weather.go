package main

type OpenWeather struct {
	Weather
	Response int `json:"cod"`
}

type Weather struct {
	ID          int           `json:"id" bson:"_id"`
	Coord       Corditations  `json:"coord" bson:"Coord"`
	WeatherData []WeatherData `json:"weather" bson:"WeatherData"`
	MainData    MainData      `json:"main" bson:"MainData"`
	Name        string        `json:"name" bson:"Name"`
}

type Corditations struct {
	Lon float64 `json:"lon" bson:"Lon"`
	Lat float64 `json:"lat" bson:"Lat"`
}

type WeatherData struct {
	ID          int    `json:"id" bson:"ID"`
	Main        string `json:"main" bson:"Main"`
	Description string `json:"description" bson:"Description"`
	Icon        string `json:"icon" bson:"Icon"`
}

type MainData struct {
	Temp      float64 `json:"temp" bson:"Temp"`
	FeelsLike float64 `json:"feels_like" bson:"FeelsLike"`
	MinTemp   float64 `json:"temp_min" bson:"MinTemp"`
	MaxTemp   float64 `json:"temp_max" bson:"MaxTemp"`
	Humidity  float64 `json:"humidity" bson:"Humidity"`
}
