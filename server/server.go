package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.RequestURI, r.Method)
		next.ServeHTTP(w, r)
	})
}

func handleMovieRequest(w http.ResponseWriter, r *http.Request) {
	movieName := mux.Vars(r)["name"]
	resp, err := http.Get(viper.GetString("omdbPrefix") + "&t=" + movieName)
	if err != nil {
		log.Println("error fetching movie detail ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	omdbMovie := &OmdbMovie{}
	parseBody(resp, omdbMovie)
	if omdbMovie.Response == viper.GetString("omdbBadResponse") {
		log.Println("request didn't yield results")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = json.NewEncoder(w).Encode(omdbMovie.Movie); err != nil {
		log.Println("error encoding result back to the caller ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func handleWeatherRequest(w http.ResponseWriter, r *http.Request) {
	cityName := mux.Vars(r)["cityName"]
	resp, err := http.Get(viper.GetString("weatherBaseUrl") + "q=" + cityName + "&" + viper.GetString("weatherApiQueryString"))
	if err != nil {
		log.Println("error fetching weather detail ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	weather := &weather{}
	parseBody(resp, weather)
	if weather.Response != viper.GetInt("weatherGoodResponse") {
		log.Println("request didn't yield results")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = json.NewEncoder(w).Encode(weather); err != nil {
		log.Println("error encoding result back to the caller ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func parseBody(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("error reading response body ", err)
		return err
	}

	if err = json.Unmarshal(body, result); err != nil {
		log.Println("error parsing response body ", err)
		return err
	}

	return nil
}

func startHttpServer() {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.HandleFunc("/api/movie/{name}", handleMovieRequest).Methods("GET")
	r.HandleFunc("/api/weather/{cityName}", handleWeatherRequest).Methods("GET")
	log.Println("listening on :", viper.GetString("port"))
	srv := &http.Server{
		Handler: r,
		Addr:    ":" + viper.GetString("port"),
	}
	log.Fatal(srv.ListenAndServe())
}
