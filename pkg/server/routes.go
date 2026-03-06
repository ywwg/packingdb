package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"sort"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ywwg/packingdb/pkg/logger"
	"github.com/ywwg/packingdb/pkg/packinglib"
)

// Handlers

// GET /api/trips - List all trips
func (s *Server) listTripsHandler(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	tripNames := make([]string, 0, len(s.nameToFile))
	for name := range s.nameToFile {
		tripNames = append(tripNames, name)
	}
	s.mu.RUnlock()

	sort.Strings(tripNames)

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

	if req.Nights <= 0 {
		s.respondError(w, "Nights must be a positive number", http.StatusBadRequest)
		return
	}

	// Hold the lock for the remainder of the handler so that the duplicate
	// check, reservation, I/O, and final mapping update are all atomic with
	// respect to other concurrent requests and the background persistence
	// goroutine.
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.nameToFile[req.Name]; exists {
		s.respondError(w, fmt.Sprintf("A trip named %q already exists", req.Name), http.StatusConflict)
		return
	}

	safeBase := sanitizeFilename(req.Name)
	filename, err := s.uniqueTripFilename(safeBase)
	if err != nil {
		s.respondError(w, fmt.Sprintf("Failed to generate filename: %v", err), http.StatusInternalServerError)
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
	s.nameToFile[req.Name] = filename
	logger.Info("Created new trip", "name", req.Name, "file", filename)

	s.respondJSON(w, map[string]interface{}{
		"success":  true,
		"filename": filepath.Base(filename),
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

	dirtyKey := name

	s.mu.Lock()
	if req.Name != nil && *req.Name != "" && *req.Name != trip.C.Name {
		// Check for name collision with an existing trip
		if _, exists := s.nameToFile[*req.Name]; exists {
			s.mu.Unlock()
			s.respondError(w, fmt.Sprintf("A trip named %q already exists", *req.Name), http.StatusConflict)
			return
		}

		oldName := trip.C.Name
		trip.C.Name = *req.Name

		// Update nameToFile mapping
		if filePath, ok := s.nameToFile[oldName]; ok {
			delete(s.nameToFile, oldName)
			s.nameToFile[*req.Name] = filePath
		}

		// Re-key trip cache
		delete(s.trips, oldName)
		s.trips[*req.Name] = trip

		// Re-key dirty trips
		if s.dirtyTrips[oldName] {
			delete(s.dirtyTrips, oldName)
		}

		dirtyKey = *req.Name
	}

	if req.Nights != nil {
		if *req.Nights <= 0 {
			s.mu.Unlock()
			s.respondError(w, "Nights must be a positive number", http.StatusBadRequest)
			return
		}
		trip.C.Nights = *req.Nights
	}
	if req.TemperatureMin != nil {
		trip.C.TemperatureMin = *req.TemperatureMin
	}
	if req.TemperatureMax != nil {
		trip.C.TemperatureMax = *req.TemperatureMax
	}

	s.dirtyTrips[dirtyKey] = true
	s.mu.Unlock()

	s.respondJSON(w, map[string]interface{}{
		"success": true,
		"name":    dirtyKey,
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

	// Hold the lock while mutating the trip so that the background persistence
	// goroutine (which also holds s.mu during SaveToFile) cannot race on the
	// trip's internal data structures.
	s.mu.Lock()
	toggleErr := trip.ToggleItemPacked(code)
	if toggleErr == nil {
		s.dirtyTrips[name] = true
	}
	s.mu.Unlock()

	if toggleErr != nil {
		s.respondError(w, fmt.Sprintf("Failed to toggle item: %v", toggleErr), http.StatusBadRequest)
		return
	}

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

	// Hold the lock while mutating the trip to prevent races with the background
	// persistence goroutine (which holds s.mu during SaveToFile).
	s.mu.Lock()
	var mutateErr error
	if trip.HasProperty(packinglib.Property(property)) {
		mutateErr = trip.RemoveProperty(property)
	} else {
		mutateErr = trip.AddProperty(property)
	}
	if mutateErr == nil {
		s.dirtyTrips[name] = true
	}
	s.mu.Unlock()

	if mutateErr != nil {
		s.respondError(w, fmt.Sprintf("Failed to toggle property: %v", mutateErr), http.StatusBadRequest)
		return
	}

	s.respondJSON(w, map[string]interface{}{
		"success": true,
	})
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
