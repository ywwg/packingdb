// packingweb is a web-based frontend for packingdb with REST API
package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ywwg/packingdb/pkg/contexts"
	"github.com/ywwg/packingdb/pkg/logger"
	"github.com/ywwg/packingdb/pkg/packinglib"
	"github.com/ywwg/packingdb/pkg/server"
)

func main() {
	// Define command-line flags
	tripsDir := flag.String("trips", "./public/trips", "Directory for trip files")
	staticDir := flag.String("static", "./static", "Directory for static files")
	port := flag.String("port", "8080", "Port to listen on")
	logLevel := flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	flag.Parse()

	// Initialize logger
	logger.Init(*logLevel)

	registry := packinglib.NewStructRegistry()
	contexts.PopulateRegistry(registry)

	// Create trips directory if it doesn't exist
	if err := os.MkdirAll(*tripsDir, 0755); err != nil {
		logger.Fatal("Failed to create trips directory", "error", err)
	}

	// Create API server (scans trip files to build name→filename mapping)
	apiServer, err := server.NewServer(registry, *tripsDir)
	if err != nil {
		logger.Fatal("Failed to create server", "error", err)
	}

	apiServer.StartBackgroundPersist()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		logger.Info("Shutting down, saving changes...")
		apiServer.Shutdown()
		os.Exit(0)
	}()

	// Serve static files
	fs := http.FileServer(http.Dir(*staticDir))
	http.Handle("/", fs)

	// API routes
	http.Handle("/api/", apiServer.Handler())

	logger.Info("Starting packingweb server", "port", *port)
	logger.Info("Visit the app", "url", "http://localhost:"+*port)
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		logger.Fatal("Server failed", "error", err)
	}
}
