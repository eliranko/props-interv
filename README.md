# Explanation
This project implements two main features:
1. Enables a user to fetch movie details from OMDB
2. Enables a user to fetch weather details from OpenWeather

# Technologies
## Backend
The backend is written in GO. It uses Mux for comfortable REST implementation and Viper for parsing the configuration file.
It also connects to MongoDB atlas to reduce 3rd party API requests.
## Frontend
The frontend is written in Angular and uses Angular material.

# Installation
Run the container eliranko/prop-interv:latest from the docker hub
``` docker run --publish 15000:80 --name pros-interv eliranko/prop-interv:latest ```
