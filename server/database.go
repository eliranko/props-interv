package main

// This file interfaces with MongoDB

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// signaling connection to MongoDB
var connectedToMongo = make(chan struct{})

// Collections
var weatherCollection *mongo.Collection
var moviesCollection *mongo.Collection

var notConnectionError = errors.New("not connected to db")

// Start a connection to MongoDB
// creates Weather & Movides collections when connection occurs
// and signals that a successful MongoDB connection occured
func startDb() {
	var client *mongo.Client
	var err error
	for {
		client, err = mongo.NewClient(options.Client().ApplyURI(viper.GetString("mongoUrl")))
		if err != nil {
			log.Panic(err)
		}

		timeout := 30 * time.Second
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		if err = client.Connect(ctx); err != nil {
			log.Println(err)
			continue
		}

		ctx, _ = context.WithTimeout(context.Background(), timeout)
		if err = client.Ping(ctx, readpref.Primary()); err != nil {
			log.Println(err)
			continue
		}
	}

	log.Println("Connected to db at ", viper.GetString("mongoUrl"))
	weatherCollection = client.Database(viper.GetString("mongoDbName")).Collection(viper.GetString("mongoWeatherCollectionName"))
	moviesCollection = client.Database(viper.GetString("mongoDbName")).Collection(viper.GetString("mongoMoviesCollectionName"))
	close(connectedToMongo)
}

// Gets a movie from the Movies collection
// This isn't efficient. The entire collections is search one every request because
// this uses ObjectId as the Index.
// This should create index on Title field for O(log n) retrieval operations
func getMovieFromDb(ctx context.Context, title string) (*Movie, error) {
	select {
	case <-connectedToMongo:
	case <-ctx.Done():
		return nil, notConnectionError
	}

	movie := &Movie{}
	err := moviesCollection.FindOne(ctx, bson.M{"Title": title}).Decode(movie)
	if err != nil {
		log.Printf("Failed reading movie %s from its collection %s", title, err)
		return nil, err
	}

	return movie, nil
}

// Gets a weather details from the Weather collection
// This isn't efficient. The entire collections is search one every request because
// this uses ObjectId as the Index.
// This should create index on Name field for O(log n) retrieval operations
func getWeatherFromDb(ctx context.Context, city string) (*Weather, error) {
	select {
	case <-connectedToMongo:
	case <-ctx.Done():
		return nil, notConnectionError
	}

	weather := &Weather{}
	err := weatherCollection.FindOne(ctx, bson.M{"Name": city}).Decode(weather)
	if err != nil {
		log.Printf("Failed reading the weather of %s from its collection %s", city, err)
		return nil, err
	}

	return weather, nil
}

// Adds a movie to the Movies collection
func persistMovie(ctx context.Context, movie *Movie) error {
	select {
	case <-connectedToMongo:
	case <-ctx.Done():
		return notConnectionError
	}

	if _, err := moviesCollection.InsertOne(ctx, movie); err != nil {
		log.Println("could not persist movie ", err)
		return err
	}

	return nil
}

// Adds weather detail to the Weather collection
func persistWeather(ctx context.Context, weather *Weather) error {
	select {
	case <-connectedToMongo:
	case <-ctx.Done():
		return notConnectionError
	}

	if _, err := weatherCollection.InsertOne(ctx, weather); err != nil {
		log.Println("could not persist movie ", err)
		return err
	}

	return nil
}
