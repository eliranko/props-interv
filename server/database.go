package main

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

var connectedToMongo = make(chan struct{})
var weatherCollection *mongo.Collection
var moviesCollection *mongo.Collection

func startDb() {
	client, err := mongo.NewClient(options.Client().ApplyURI(viper.GetString("mongoUrl")))
	if err != nil {
		log.Panic(err)
	}

	timeout := 30 * time.Second
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	if err = client.Connect(ctx); err != nil {
		log.Println(err)
		return
	}

	ctx, _ = context.WithTimeout(context.Background(), timeout)
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Println(err)
		return
	}

	log.Println("Connected to db at ", viper.GetString("mongoUrl"))
	weatherCollection = client.Database(viper.GetString("mongoDbName")).Collection(viper.GetString("mongoWeatherCollectionName"))
	moviesCollection = client.Database(viper.GetString("mongoDbName")).Collection(viper.GetString("mongoMoviesCollectionName"))
	close(connectedToMongo)
}

func getMovieFromDb(ctx context.Context, title string) (*Movie, error) {
	select {
	case <-connectedToMongo:
	case <-ctx.Done():
		return nil, errors.New("not connected to db")
	}

	movie := &Movie{}
	err := moviesCollection.FindOne(ctx, bson.M{"Title": title}).Decode(movie)
	if err != nil {
		log.Printf("Failed reading movie %s from its collection %s", title, err)
		return nil, err
	}

	return movie, nil
}

func getWeatherFromDb(ctx context.Context, city string) (*Weather, error) {
	select {
	case <-connectedToMongo:
	case <-ctx.Done():
		return nil, errors.New("not connected to db")
	}

	weather := &Weather{}
	err := weatherCollection.FindOne(ctx, bson.M{"Name": city}).Decode(weather)
	if err != nil {
		log.Printf("Failed reading the weather of %s from its collection %s", city, err)
		return nil, err
	}

	return weather, nil
}

func persistMovie(ctx context.Context, movie *Movie) error {
	select {
	case <-connectedToMongo:
	case <-ctx.Done():
		return errors.New("not connected to db")
	}

	if _, err := moviesCollection.InsertOne(ctx, movie); err != nil {
		log.Println("could not persist movie ", err)
		return err
	}

	return nil
}

func persistWeather(ctx context.Context, weather *Weather) error {
	select {
	case <-connectedToMongo:
	case <-ctx.Done():
		return errors.New("not connected to db")
	}

	if _, err := weatherCollection.InsertOne(ctx, weather); err != nil {
		log.Println("could not persist movie ", err)
		return err
	}

	return nil
}
