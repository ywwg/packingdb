// packingweb is a web-based frontend for packingdb with REST API
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ywwg/packingdb/pkg/contexts"
	"github.com/ywwg/packingdb/pkg/packinglib"
	"github.com/ywwg/packingdb/pkg/server"
)

func main() {
	// Define command-line flags
	tripsDir := flag.String("trips", "./public/trips", "Directory for trip files")
	staticDir := flag.String("static", "./static", "Directory for static files")
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	registry := packinglib.NewStructRegistry()
	contexts.PopulateRegistry(registry)

	// Create trips directory if it doesn't exist
	if err := os.MkdirAll(*tripsDir, 0755); err != nil {
		log.Fatal(err)
	}

	// Create API server
	apiServer := server.NewServer(registry, *tripsDir)
	apiServer.StartBackgroundPersist()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down, saving changes...")
		apiServer.Shutdown()
		os.Exit(0)
	}()

	// Serve static files
	fs := http.FileServer(http.Dir(*staticDir))
	http.Handle("/", fs)

	// API routes
	http.Handle("/api/", apiServer.Handler())

	log.Printf("Starting packingweb server on port %s", *port)
	log.Printf("Visit http://localhost:%s to use the app", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
