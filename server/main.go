package main

import (
	"log"

	"github.com/spf13/viper"
)

// loads up the fields from the configuration file into memory
func init() {
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Panicln(err)
	}
}

// starts the DB connection and the HTTP server
func main() {
	go startDb()
	startHTTPServer()
}
