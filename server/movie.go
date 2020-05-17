package main

type OmdbMovie struct {
	Movie
	Response string `json:"Response"`
}

type Movie struct {
	Title    string `json:"Title"`
	Year     string `json:"Year"`
	Plot     string `json:"Plot"`
	Language string `json:"Language"`
	Poster   string `json:"Poster"`
	Rating   string `json:"imdbRating"`
}
