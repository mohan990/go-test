package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type messageResponse struct {
	Message string `json:"message"`
	Time    string `json:"time"`
}

type healthResponse struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

// RootHandler handles requests to the root endpoint
func RootHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, messageResponse{
		Message: "Go webserver is running",
		Time:    time.Now().UTC().Format(time.RFC3339),
	})
}

// HelloHandler handles requests to the hello endpoint
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "world"
	}
	writeJSON(w, http.StatusOK, messageResponse{
		Message: "Hello, " + name + "!",
		Time:    time.Now().UTC().Format(time.RFC3339),
	})
}

// HealthzHandler handles health check requests
func HealthzHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, healthResponse{
		Status: "ok",
		Time:   time.Now().UTC().Format(time.RFC3339),
	})
}

// ReadinessHandler handles readiness probe requests
func ReadinessHandler(w http.ResponseWriter, _ *http.Request) {
	// Add any readiness checks here (e.g., database connection, external services)
	// For now, we'll always return ready
	writeJSON(w, http.StatusOK, healthResponse{
		Status: "ready",
		Time:   time.Now().UTC().Format(time.RFC3339),
	})
}

// LivenessHandler handles liveness probe requests
func LivenessHandler(w http.ResponseWriter, _ *http.Request) {
	// Liveness check - is the application alive and not deadlocked
	writeJSON(w, http.StatusOK, healthResponse{
		Status: "alive",
		Time:   time.Now().UTC().Format(time.RFC3339),
	})
}
