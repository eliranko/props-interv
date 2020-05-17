package main

// This file handles the REST API through HTTP
// It uses in-memory caching for incoming request

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

// Should use Reddis instead of a HashMap
// This implementation isn't safe nor efficient:
// 1. no capacity - this can use the entire virtual memory
// 2. No cleaning policy (like LRU)
var httpMoviesCache = make(map[string]*Movie)
var httpWeatherCache = make(map[string]*Weather)

// Timeout for DB or HTTP requests
const httpTimeout = 5 * time.Second

var httpClient = http.Client{
	Timeout: httpTimeout,
}

// logs incoming requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.RequestURI, r.Method)
		next.ServeHTTP(w, r)
	})
}

// handles GET requests for some movie details
// the request is fetched from cache, if it exists
// otherwise, it will be fetched from the DB
// if both failed (not available or don't have the record), the request
// will be fecthed through HTTP from OMDB
func handleMovieRequest(w http.ResponseWriter, r *http.Request) {
	// Set context for DB request, if needed
	ctx, _ := context.WithTimeout(context.Background(), httpTimeout)
	// Changed movie name to consistent notation
	movieName := consistentNotation(mux.Vars(r)["name"])
	var omdbMovie *OmdbMovie

	if movie, ok := httpMoviesCache[movieName]; ok { // load from cache
		omdbMovie = &OmdbMovie{Movie: *movie}
		log.Printf("got %s from cache", movieName)

	} else if movie, err := getMovieFromDb(ctx, movieName); err == nil { // load from DB
		omdbMovie = &OmdbMovie{Movie: *movie}
		log.Printf("got %s from db", movieName)

	} else { // Fetch from OMDB
		log.Printf("coudln't fetch %s from cache or db, sending HTTP request...", movieName)
		resp, err := httpClient.Get(viper.GetString("omdbPrefix") + viper.GetString("omdbMovieNameQueryString") + movieName)
		if err != nil {
			log.Printf("error %s fetching movie %s from OMDB", err, movieName)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		omdbMovie = &OmdbMovie{}
		parseBody(resp, omdbMovie)
		// Check if received valid response
		if omdbMovie.Response == viper.GetString("omdbBadResponse") {
			log.Println("request didn't yield any result")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		omdbMovie.Title = consistentNotation(omdbMovie.Title) // consistent notation

		// Save data for follow-up request
		httpMoviesCache[movieName] = &omdbMovie.Movie
		go func() {
			ctx, _ := context.WithTimeout(context.Background(), httpTimeout)
			_ = persistMovie(ctx, &omdbMovie.Movie) // ignore error
		}()
	}

	// Send result to user
	if err := json.NewEncoder(w).Encode(omdbMovie.Movie); err != nil {
		log.Println("error encoding result back to the caller ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// handles GET requests for some city's weather details
// the request is fetched from cache, if it exists
// otherwise, it will be fetched from the DB
// if both failed (not available or don't have the record), the request
// will be fecthed through HTTP from OpenWeather
func handleWeatherRequest(w http.ResponseWriter, r *http.Request) {
	// Set context for DB request, if needed
	ctx, _ := context.WithTimeout(context.Background(), httpTimeout)
	// Changed city name to consistent notation
	cityName := consistentNotation(mux.Vars(r)["cityName"])
	var openWeather *OpenWeather

	if weather, ok := httpWeatherCache[cityName]; ok { // load from cache
		openWeather = &OpenWeather{Weather: *weather}
		log.Printf("got %s from cache", cityName)

	} else if weather, err := getWeatherFromDb(ctx, cityName); err == nil { // load from DB
		openWeather = &OpenWeather{Weather: *weather}
		log.Printf("got %s from db", cityName)

	} else { // Fetch from OpenWeather
		log.Printf("coudln't fetch %s from cache or db, sending HTTP request...", cityName)
		resp, err := httpClient.Get(viper.GetString("weatherBaseUrl") + "q=" + cityName + "&" + viper.GetString("weatherApiQueryString"))
		if err != nil {
			log.Printf("error %s fetching weather details of %s", err, cityName)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		openWeather = &OpenWeather{}
		parseBody(resp, openWeather)
		// Check if received valid result
		if openWeather.Response != viper.GetInt("weatherGoodResponse") {
			log.Println("request didn't yield any result")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		openWeather.Name = consistentNotation(openWeather.Name) // consistent notation

		// Save data for follow-up request
		httpWeatherCache[cityName] = &openWeather.Weather
		go func() {
			ctx, _ := context.WithTimeout(context.Background(), httpTimeout)
			_ = persistWeather(ctx, &openWeather.Weather) // ignore error
		}()
	}

	// Send result to user
	if err := json.NewEncoder(w).Encode(openWeather.Weather); err != nil {
		log.Println("error encoding result back to the caller ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func consistentNotation(name string) string {
	return strings.ToTitle(name)
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

func startHTTPServer() {
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
