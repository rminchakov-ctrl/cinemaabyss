package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/cinemaabyss/src/microservices/proxy/models"
	"github.com/cinemaabyss/src/microservices/proxy/services"
)

var moviesService *services.MovieService

func main() {
	moviesMigration := getEnvBool("GRADUAL_MIGRATION", true)
	moviesMigrationPercent := getEnvInt("MOVIES_MIGRATION_PERCENT", 50)
	log.Printf("Migranion movies to new service? %t Percent: %d\n", moviesMigration, moviesMigrationPercent)
	if !moviesMigration {
		moviesMigrationPercent = 0
	}

	monolithURL := getEnv("MONOLITH_URL", "http://monolith:8080")
	moviesAPIURL := getEnv("MOVIES_SERVICE_URL", "http://movies-service:8081")
	moviesService = services.NewMoviesService(monolithURL, moviesAPIURL, moviesMigrationPercent)
	log.Printf("Movies service initialized with service URL: %s, monolith URL %s\n", moviesAPIURL, moviesAPIURL)

	/*
		eventAPIURL := getEnv("EVENTS_SERVICE_URL", "http://events-service:8082")
		eventService := services.NewEventService(eventAPIURL)
		log.Printf("Event service initialized with API URL: %s\n", eventAPIURL)
	*/
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/movies", handleMovies)

	port := getEnv("PORT", "8000")
	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvBool gets an environment variable or returns a default value
func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	if retVal, err := strconv.ParseBool(value); err != nil {
		return retVal
	}
	return defaultValue
}

// getEnvInt gets an environment variable or returns a default value
func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	if retVal, err := strconv.Atoi(value); err != nil {
		return retVal
	}
	return defaultValue
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"status": true})
}

func handleMovies(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	switch r.Method {
	case "GET":
		if r.URL.Query().Get("id") != "" {
			id := r.URL.Query().Get("id")
			movie, err := moviesService.GetMovieByID(ctx, id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(movie)
		} else {
			movies, err := moviesService.GetMovies(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(movies)
		}
	case "POST":
		var movie models.Movie
		if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		movieNew, err := moviesService.CreateMovie(ctx, &movie)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(movieNew)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
