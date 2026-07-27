package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type frontendResponse struct {
	Service         string         `json:"service"`
	Message         string         `json:"message"`
	BackendResponse map[string]any `json:"backend_response"`
}

func main() {
	port := getEnv("PORT", "8081")
	backendURL := strings.TrimRight(getEnv("BACKEND_URL", "http://localhost:8080"), "/")
	client := &http.Client{Timeout: 5 * time.Second}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"status":  "ok",
			"service": "frontend",
		})
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		resp, err := client.Get(backendURL + "/api/message")
		if err != nil {
			writeJSON(w, http.StatusBadGateway, map[string]string{
				"service": "frontend",
				"error":   fmt.Sprintf("backend request failed: %v", err),
			})
			return
		}
		defer resp.Body.Close()

		var backendPayload map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&backendPayload); err != nil {
			writeJSON(w, http.StatusBadGateway, map[string]string{
				"service": "frontend",
				"error":   fmt.Sprintf("invalid backend response: %v", err),
			})
			return
		}

		writeJSON(w, http.StatusOK, frontendResponse{
			Service:         "frontend",
			Message:         "Frontend successfully contacted the backend",
			BackendResponse: backendPayload,
		})
	})

	log.Printf("frontend listening on :%s and using backend %s", port, backendURL)
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
