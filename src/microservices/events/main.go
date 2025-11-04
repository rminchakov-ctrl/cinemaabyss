package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cinemaabyss/src/microservices/events/models"
	"github.com/cinemaabyss/src/microservices/events/services"
)

var eventService *services.EventService

func main() {
	ctx := context.Background()

	// Set up HTTP routes
	http.HandleFunc("/api/events", handleEvents)
	http.HandleFunc("/api/events/health", handleHealth)

	kafkaBrokers := getEnv("KAFKA_BROKERS", "kafka:9092")
	eventService = services.NewEventService(kafkaBrokers)
	defer eventService.Close()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			select {
			case <-sigchan:
				log.Println("Shutting down gracefully...")
				return

			default:
				msg, err := eventService.Read(ctx)
				if err != nil {
					log.Printf("Error reading message: %v\n", err)
					continue
				}
				log.Printf("Received message: key=%s value=%s", msg.Key, msg.Value)
			}
		}
	}()

	// Start server
	port := getEnv("PORT", "8082")
	log.Printf("Starting events microservice on port %s", port)
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

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"status": true})
}

func handleEvents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var event models.Event
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ok := true
		err := eventService.Write(context.Background(), &event)
		if err != nil {
			ok = false
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"status": ok})
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
