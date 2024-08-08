package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/codebyharshit/real-time-analytics/internal/core/entities"
	"github.com/codebyharshit/real-time-analytics/internal/infrastructure/di"
	"github.com/rs/cors"
)

func main() {
	container, err := di.NewContainer()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/trade", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var trade entities.Trade
		if err := json.NewDecoder(r.Body).Decode(&trade); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := container.TraderService.ExecuteTrade(trade); err != nil {
			http.Error(w, "Failed to execute trade", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/portfolio", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		portfolio, err := container.TraderService.GetPortfolio()
		if err != nil {
			http.Error(w, "Failed to get portfolio", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(portfolio); err != nil {
			http.Error(w, "Failed to encode portfolio", http.StatusInternalServerError)
			return
		}
	})

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	handler := corsHandler.Handler(http.DefaultServeMux)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
