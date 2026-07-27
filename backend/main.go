package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type response map[string]string

func main() {
	port := getEnv("PORT", "8080")

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, response{
			"status":  "ok",
			"service": "backend",
		})
	})
	mux.HandleFunc("/api/message", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, response{
			"service": "backend",
			"message": "Hello from the backend microservice",
		})
	})

	log.Printf("backend listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
