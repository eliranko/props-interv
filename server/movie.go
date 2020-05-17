package main

type OmdbMovie struct {
	Movie
	Response string `json:"Response"`
}

type Movie struct {
	Title    string `json:"Title" bson:"Title"`
	Year     string `json:"Year" bson:"Year"`
	Plot     string `json:"Plot" bson:"Plot"`
	Language string `json:"Language" bson:"Language"`
	Poster   string `json:"Poster" bson:"Poster"`
	Rating   string `json:"imdbRating" bson:"Rating"`
}
