package main

import (
	"log"

	"github.com/spf13/viper"
)

func init() {
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Panicln(err)
	}
}

func main() {
	startHttpServer()
}
