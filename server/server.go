package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

var httpMoviesCache = make(map[string]*Movie)
var httpWeatherCache = make(map[string]*Weather)

const httpTimeout = 5 * time.Second

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.RequestURI, r.Method)
		next.ServeHTTP(w, r)
	})
}

func handleMovieRequest(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), httpTimeout)
	movieName := strings.ToLower(mux.Vars(r)["name"])
	var omdbMovie *OmdbMovie

	if movie, ok := httpMoviesCache[movieName]; ok { // load from cache
		omdbMovie = &OmdbMovie{Movie: *movie}
		log.Printf("got %s from cache", movieName)
	} else if movie, err := getMovieFromDb(ctx, movieName); err == nil { // load from DB
		omdbMovie = &OmdbMovie{Movie: *movie}
		log.Printf("got %s from db", movieName)
	} else { // Fetch from OMDB
		resp, err := http.Get(viper.GetString("omdbPrefix") + "&t=" + movieName)
		if err != nil {
			log.Println("error fetching movie detail ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		omdbMovie = &OmdbMovie{}
		parseBody(resp, omdbMovie)
		if omdbMovie.Response == viper.GetString("omdbBadResponse") {
			log.Println("request didn't yield results")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		httpMoviesCache[movieName] = &omdbMovie.Movie
		ctx, _ := context.WithTimeout(context.Background(), httpTimeout)
		_ = persistMovie(ctx, &omdbMovie.Movie) // ignore error
	}

	if err := json.NewEncoder(w).Encode(omdbMovie.Movie); err != nil {
		log.Println("error encoding result back to the caller ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func handleWeatherRequest(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), httpTimeout)
	cityName := strings.ToLower(mux.Vars(r)["cityName"])
	var openWeather *OpenWeather

	if weather, ok := httpWeatherCache[cityName]; ok { // load from cache
		openWeather = &OpenWeather{Weather: *weather}
		log.Printf("got %s from cache", cityName)
	} else if weather, err := getWeatherFromDb(ctx, cityName); err == nil { // load from DB
		openWeather = &OpenWeather{Weather: *weather}
		log.Printf("got %s from db", cityName)
	} else { // Fetch from OMDB
		resp, err := http.Get(viper.GetString("weatherBaseUrl") + "q=" + cityName + "&" + viper.GetString("weatherApiQueryString"))
		if err != nil {
			log.Println("error fetching movie detail ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		openWeather = &OpenWeather{}
		parseBody(resp, openWeather)
		if openWeather.Response != viper.GetInt("weatherGoodResponse") {
			log.Println("request didn't yield results")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		httpWeatherCache[cityName] = &openWeather.Weather
		ctx, _ := context.WithTimeout(context.Background(), httpTimeout)
		_ = persistWeather(ctx, &openWeather.Weather) // ignore error
	}

	if err := json.NewEncoder(w).Encode(openWeather.Weather); err != nil {
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
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./dist/")))

	log.Println("listening on :", viper.GetString("port"))
	srv := &http.Server{
		Handler: r,
		Addr:    ":" + viper.GetString("port"),
	}
	log.Fatal(srv.ListenAndServe())
}
