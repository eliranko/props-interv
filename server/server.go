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
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("error reading omdb response body ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	omdbMovie := &OmdbMovie{}
	if err = json.Unmarshal(body, omdbMovie); err != nil {
		log.Println("error parsing omdb body ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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

func startHttpServer() {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.HandleFunc("/api/movie/{name}", handleMovieRequest).Methods("GET")
	log.Println("listening on :", viper.GetString("port"))
	srv := &http.Server{
		Handler: r,
		Addr:    ":" + viper.GetString("port"),
	}
	log.Fatal(srv.ListenAndServe())
}
