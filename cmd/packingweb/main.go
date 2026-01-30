// packingweb is a web-based frontend for packingdb with REST API
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ywwg/packingdb/pkg/api"
	"github.com/ywwg/packingdb/pkg/contexts"
	"github.com/ywwg/packingdb/pkg/packinglib"
)

func main() {
	// Initialize registry
	registry := packinglib.NewStructRegistry()
	contexts.PopulateRegistry(registry)

	// Determine trips directory
	// Can be overridden with TRIPS_DIR environment variable for testing
	tripsDir := os.Getenv("TRIPS_DIR")
	if tripsDir == "" {
		// Try to use ../../public/trips (from cmd/packingweb to repo root)
		// If that path's parent exists, use it; otherwise use public/trips
		tripsDir = "../../public/trips"
		if _, err := os.Stat("../../static"); os.IsNotExist(err) {
			// Not running from cmd/packingweb, must be running from repo root
			tripsDir = "public/trips"
		}
	}

	// Create trips directory if it doesn't exist
	if err := os.MkdirAll(tripsDir, 0755); err != nil {
		log.Fatal(err)
	}

	// Determine static files directory
	// Can be overridden with STATIC_DIR environment variable for testing
	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		// If running from cmd/packingweb, static is at ../../static
		// If running from repo root, static is at ./static
		staticDir = "../../static"
		if _, err := os.Stat(staticDir); os.IsNotExist(err) {
			staticDir = "static"
		}
	}

	// Create API server
	apiServer := api.NewServer(registry, tripsDir)

	// Serve static files
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/", fs)

	// API routes
	http.Handle("/api/", apiServer.Handler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting packingweb server on port %s", port)
	log.Printf("Visit http://localhost:%s to use the app", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
