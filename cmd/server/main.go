package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/your-username/go-webserver-jenkins-gke/internal/handler"
)

func newMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Application routes
	mux.HandleFunc("/", handler.RootHandler)
	mux.HandleFunc("/hello", handler.HelloHandler)

	// Health check routes (separated for better observability)
	mux.HandleFunc("/healthz", handler.HealthzHandler)
	mux.HandleFunc("/readyz", handler.ReadinessHandler)
	mux.HandleFunc("/livez", handler.LivenessHandler)

	return mux
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           newMux(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("server listening on :%s", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
