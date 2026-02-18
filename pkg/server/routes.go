package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ywwg/packingdb/pkg/logger"
	"github.com/ywwg/packingdb/pkg/packinglib"
)

// Server holds the state for the API server
type Server struct {
	Registry packinglib.Registry
	TripsDir string

	mu         sync.RWMutex
	trips      map[string]*packinglib.Trip
	dirtyTrips map[string]bool // tracks which trips have unsaved changes

	// For background persistence
	persistInterval time.Duration
	stopPersist     chan struct{}
	persistDone     chan struct{}
}

// NewServer creates a new API server instance
func NewServer(registry packinglib.Registry, tripsDir string) *Server {
	return &Server{
		Registry:        registry,
		TripsDir:        tripsDir,
		trips:           make(map[string]*packinglib.Trip),
		dirtyTrips:      make(map[string]bool),
		persistInterval: 30 * time.Second,
		stopPersist:     make(chan struct{}),
		persistDone:     make(chan struct{}),
	}
}

// Response types
type ErrorResponse struct {
	Error string `json:"error"`
}

type TripInfo struct {
	Name           string   `json:"name"`
	Nights         int      `json:"nights"`
	TemperatureMin int      `json:"temperatureMin"`
	TemperatureMax int      `json:"temperatureMax"`
	Properties     []string `json:"properties"`
}

type ItemResponse struct {
	Code     string  `json:"code"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Packed   bool    `json:"packed"`
	Count    float64 `json:"count"`
}

type CategoryResponse struct {
	Name  string         `json:"name"`
	Items []ItemResponse `json:"items"`
}

type PropertyResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

// Helper functions
func (s *Server) respondError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func (s *Server) respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// Handlers

// GET /api/trips - List all trips
func (s *Server) listTripsHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(s.TripsDir)
	if err != nil {
		s.respondError(w, "Failed to read trips directory", http.StatusInternalServerError)
		return
	}

	tripNames := []string{}
	for _, file := range files {
		if !file.IsDir() && (strings.HasSuffix(file.Name(), ".csv") ||
			strings.HasSuffix(file.Name(), ".yml") ||
			strings.HasSuffix(file.Name(), ".yaml")) {
			tripNames = append(tripNames, file.Name())
		}
	}

	s.respondJSON(w, map[string]interface{}{
		"trips": tripNames,
	})
}

// POST /api/trips - Create a new trip
func (s *Server) createTripHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name           string   `json:"name"`
		Nights         int      `json:"nights"`
		TemperatureMin int      `json:"temperatureMin"`
		TemperatureMax int      `json:"temperatureMax"`
		Properties     []string `json:"properties"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		s.respondError(w, "Trip name is required", http.StatusBadRequest)
		return
	}

	filename := filepath.Join(s.TripsDir, req.Name+".yml")

	// Check if trip already exists
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		s.respondError(w, "Trip already exists", http.StatusConflict)
		return
	}

	// Create context with properties
	context, err := packinglib.NewContext(s.Registry, req.Name, req.Nights, req.TemperatureMin, req.TemperatureMax, req.Properties)
	if err != nil {
		s.respondError(w, fmt.Sprintf("Failed to create context: %v", err), http.StatusInternalServerError)
		return
	}

	trip, err := packinglib.NewTripFromCustomContext(s.Registry, context)
	if err != nil {
		s.respondError(w, fmt.Sprintf("Failed to create trip: %v", err), http.StatusInternalServerError)
		return
	}

	if err := trip.SaveToFile(filename); err != nil {
		s.respondError(w, fmt.Sprintf("Failed to save trip: %v", err), http.StatusInternalServerError)
		return
	}

	s.trips[req.Name] = trip

	s.respondJSON(w, map[string]interface{}{
		"success":  true,
		"filename": filename,
	})
}

// GET /api/trips/{name} - Get trip details
func (s *Server) getTripHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	trip, err := s.loadTrip(name)
	if err != nil {
		s.respondError(w, fmt.Sprintf("Failed to load trip: %v", err), http.StatusNotFound)
		return
	}

	properties := []string{}
	for p, val := range trip.C.Properties {
		if val {
			properties = append(properties, string(p))
		}
	}
	sort.Strings(properties)

	s.respondJSON(w, TripInfo{
		Name:           trip.C.Name,
		Nights:         trip.C.Nights,
		TemperatureMin: trip.C.TemperatureMin,
		TemperatureMax: trip.C.TemperatureMax,
		Properties:     properties,
	})
}

// PUT /api/trips/{name} - Update trip properties
func (s *Server) updateTripHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	trip, err := s.loadTrip(name)
	if err != nil {
		s.respondError(w, fmt.Sprintf("Failed to load trip: %v", err), http.StatusNotFound)
		return
	}

	var req struct {
		Name           *string `json:"name"`
		Nights         *int    `json:"nights"`
		TemperatureMin *int    `json:"temperatureMin"`
		TemperatureMax *int    `json:"temperatureMax"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name != nil && *req.Name != "" {
		trip.C.Name = *req.Name
	}
	if req.Nights != nil {
		trip.C.Nights = *req.Nights
	}
	if req.TemperatureMin != nil {
		trip.C.TemperatureMin = *req.TemperatureMin
	}
	if req.TemperatureMax != nil {
		trip.C.TemperatureMax = *req.TemperatureMax
	}

	// Mark as dirty instead of saving immediately
	s.markDirty(name)

	s.respondJSON(w, map[string]interface{}{
		"success": true,
	})
}

// GET /api/trips/{name}/items - Get packing items for a trip
func (s *Server) getItemsHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	trip, err := s.loadTrip(name)
	if err != nil {
		s.respondError(w, fmt.Sprintf("Failed to load trip: %v", err), http.StatusNotFound)
		return
	}

	categories := []CategoryResponse{}
	currentCategory := ""
	var categoryItems []ItemResponse

	items := trip.PackingMenuItems(make(map[packinglib.Category]bool), false)
	for _, menuItem := range items {
		if menuItem.Type == packinglib.MenuCategory {
			// Save previous category if it has items
			if currentCategory != "" && len(categoryItems) > 0 {
				categories = append(categories, CategoryResponse{
					Name:  currentCategory,
					Items: categoryItems,
				})
			}
			// Start new category
			currentCategory = menuItem.Code
			categoryItems = []ItemResponse{}
		} else if menuItem.Type == packinglib.MenuPackable {
			// Get the item from the trip's internal structures
			if item, ok := trip.GetItemByCode(menuItem.Code); ok {
				categoryItems = append(categoryItems, ItemResponse{
					Code:     menuItem.Code,
					Name:     item.String(),
					Category: currentCategory,
					Packed:   item.Packed(),
					Count:    item.Count(),
				})
			}
		}
	}
	// Don't forget the last category
	if currentCategory != "" && len(categoryItems) > 0 {
		categories = append(categories, CategoryResponse{
			Name:  currentCategory,
			Items: categoryItems,
		})
	}

	s.respondJSON(w, map[string]interface{}{
		"categories": categories,
	})
}

// POST /api/trips/{name}/items/{code}/toggle - Toggle item packed status
func (s *Server) toggleItemHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	code := chi.URLParam(r, "code")

	trip, err := s.loadTrip(name)
	if err != nil {
		s.respondError(w, fmt.Sprintf("Failed to load trip: %v", err), http.StatusNotFound)
		return
	}

	if err := trip.ToggleItemPacked(code); err != nil {
		s.respondError(w, fmt.Sprintf("Failed to toggle item: %v", err), http.StatusBadRequest)
		return
	}

	// Mark as dirty instead of saving immediately
	s.markDirty(name)

	s.respondJSON(w, map[string]interface{}{
		"success": true,
	})
}

// GET /api/properties - Get all available properties
func (s *Server) getPropertiesHandler(w http.ResponseWriter, r *http.Request) {
	properties := s.Registry.ListProperties()
	propList := []PropertyResponse{}

	for _, p := range properties {
		propList = append(propList, PropertyResponse{
			Name:        string(p),
			Description: s.Registry.GetDescription(p),
			Active:      false,
		})
	}

	s.respondJSON(w, map[string]interface{}{
		"properties": propList,
	})
}

// GET /api/trips/{name}/properties - Get trip properties
func (s *Server) getTripPropertiesHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	trip, err := s.loadTrip(name)
	if err != nil {
		s.respondError(w, fmt.Sprintf("Failed to load trip: %v", err), http.StatusNotFound)
		return
	}

	properties := s.Registry.ListProperties()
	propList := []PropertyResponse{}

	for _, p := range properties {
		propList = append(propList, PropertyResponse{
			Name:        string(p),
			Description: s.Registry.GetDescription(p),
			Active:      trip.HasProperty(p),
		})
	}

	s.respondJSON(w, map[string]interface{}{
		"properties": propList,
	})
}

// POST /api/trips/{name}/properties/{property}/toggle - Toggle property
func (s *Server) togglePropertyHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	property := chi.URLParam(r, "property")

	trip, err := s.loadTrip(name)
	if err != nil {
		s.respondError(w, fmt.Sprintf("Failed to load trip: %v", err), http.StatusNotFound)
		return
	}

	if trip.HasProperty(packinglib.Property(property)) {
		if err := trip.RemoveProperty(property); err != nil {
			s.respondError(w, fmt.Sprintf("Failed to remove property: %v", err), http.StatusBadRequest)
			return
		}
	} else {
		if err := trip.AddProperty(property); err != nil {
			s.respondError(w, fmt.Sprintf("Failed to add property: %v", err), http.StatusBadRequest)
			return
		}
	}

	// Mark as dirty instead of saving immediately
	s.markDirty(name)

	s.respondJSON(w, map[string]interface{}{
		"success": true,
	})
}

// Helper functions
func (s *Server) loadTrip(name string) (*packinglib.Trip, error) {
	// Check cache first
	if trip, ok := s.trips[name]; ok {
		return trip, nil
	}

	// Find the trip file
	filename := s.findTripFile(name)
	if filename == "" {
		return nil, fmt.Errorf("trip not found: %s", name)
	}

	trip, err := packinglib.LoadFromFile(s.Registry, 0, filename)
	if err != nil {
		return nil, err
	}

	s.trips[name] = trip
	return trip, nil
}

func (s *Server) findTripFile(name string) string {
	// Try different extensions
	for _, ext := range []string{".yml", ".yaml", ".csv"} {
		filename := filepath.Join(s.TripsDir, name+ext)
		if _, err := os.Stat(filename); !os.IsNotExist(err) {
			return filename
		}
	}
	return ""
}

// Handler returns the HTTP handler for API routes
func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()

	// Request logging middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer to capture status code
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)
			logger.Info("HTTP request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", wrapped.statusCode,
				"duration", duration,
			)
		})
	})

	// CORS middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// Register routes
	r.Get("/api/trips", s.listTripsHandler)
	r.Post("/api/trips", s.createTripHandler)
	r.Get("/api/properties", s.getPropertiesHandler)
	r.Get("/api/trips/{name}", s.getTripHandler)
	r.Put("/api/trips/{name}/update", s.updateTripHandler)
	r.Get("/api/trips/{name}/items", s.getItemsHandler)
	r.Post("/api/trips/{name}/items/{code}/toggle", s.toggleItemHandler)
	r.Get("/api/trips/{name}/properties", s.getTripPropertiesHandler)
	r.Post("/api/trips/{name}/properties/{property}/toggle", s.togglePropertyHandler)

	return r
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
