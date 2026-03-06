package server

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/ywwg/packingdb/pkg/logger"
	"github.com/ywwg/packingdb/pkg/packinglib"
	yaml "gopkg.in/yaml.v3"
)

// Server holds the state for the API server
type Server struct {
	Registry packinglib.Registry
	TripsDir string

	mu         sync.RWMutex
	trips      map[string]*packinglib.Trip
	dirtyTrips map[string]bool   // tracks which trips have unsaved changes
	nameToFile map[string]string // trip name → full file path

	// For background persistence
	persistInterval time.Duration
	stopPersist     chan struct{}
	persistDone     chan struct{}
}

// NewServer creates a new API server instance. It scans the trips directory
// to build a name→filename mapping so that API lookups use the human-readable
// trip name stored inside each file rather than the filename on disk.
func NewServer(registry packinglib.Registry, tripsDir string) (*Server, error) {
	s := &Server{
		Registry:        registry,
		TripsDir:        tripsDir,
		trips:           make(map[string]*packinglib.Trip),
		dirtyTrips:      make(map[string]bool),
		nameToFile:      make(map[string]string),
		persistInterval: 30 * time.Second,
		stopPersist:     make(chan struct{}),
		persistDone:     make(chan struct{}),
	}
	if err := s.scanTrips(); err != nil {
		return nil, err
	}
	return s, nil
}

// scanTrips reads all trip files in the trips directory and builds a mapping
// from trip name (stored inside the file) to file path on disk.
// Called during construction before any concurrent access, so no mutex needed.
func (s *Server) scanTrips() error {
	files, err := os.ReadDir(s.TripsDir)
	if err != nil {
		return fmt.Errorf("failed to read trips directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		ext := filepath.Ext(file.Name())
		if ext != ".csv" && ext != ".yml" && ext != ".yaml" {
			continue
		}

		filePath := filepath.Join(s.TripsDir, file.Name())
		name, err := extractTripName(filePath)
		if err != nil {
			logger.Warn("Failed to extract name from trip file", "file", file.Name(), "error", err)
			continue
		}

		if existing, ok := s.nameToFile[name]; ok {
			return fmt.Errorf("duplicate trip name %q in files %s and %s", name, filepath.Base(existing), file.Name())
		}

		s.nameToFile[name] = filePath
		logger.Debug("Scanned trip", "name", name, "file", file.Name())
	}

	logger.Info("Trip scan complete", "count", len(s.nameToFile))
	return nil
}

// extractTripName reads just enough of a trip file to extract the trip name
// without fully loading and constructing the trip.
func extractTripName(filePath string) (string, error) {
	switch filepath.Ext(filePath) {
	case ".yml", ".yaml":
		return extractYAMLTripName(filePath)
	case ".csv":
		return extractCSVTripName(filePath)
	default:
		return "", fmt.Errorf("unsupported file type: %s", filepath.Ext(filePath))
	}
}

func extractYAMLTripName(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var partial struct {
		Name string `yaml:"name"`
	}
	if err := yaml.NewDecoder(f).Decode(&partial); err != nil {
		return "", err
	}
	if partial.Name == "" {
		return "", fmt.Errorf("no name field in YAML file")
	}
	return partial.Name, nil
}

func extractCSVTripName(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		return "", fmt.Errorf("empty file")
	}

	line := scanner.Text()
	toks := strings.Split(line, ",")

	if toks[0] == "V2" && len(toks) >= 5 {
		return toks[4], nil
	}

	if len(toks) >= 2 {
		return toks[1], nil
	}

	return "", fmt.Errorf("unrecognized CSV format")
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

// loadTrip returns a cached trip or loads it from disk on first access.
func (s *Server) loadTrip(name string) (*packinglib.Trip, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if trip, ok := s.trips[name]; ok {
		return trip, nil
	}

	filename, ok := s.nameToFile[name]
	if !ok {
		return nil, fmt.Errorf("trip not found: %s", name)
	}

	trip, err := packinglib.LoadFromFile(s.Registry, 0, filename)
	if err != nil {
		return nil, err
	}

	s.trips[name] = trip
	return trip, nil
}

// respondError sends a JSON error response with the given status code.
func (s *Server) respondError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// respondJSON sends a JSON success response.
func (s *Server) respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// sanitizeFilename removes all characters that are not alphanumeric, hyphens, or
// underscores, collapses runs of separators, and trims leading/trailing
// separators. It also strips any path components to prevent traversal attacks.
var reUnsafe = regexp.MustCompile(`[^a-zA-Z0-9_-]+`)

func sanitizeFilename(name string) string {
	// Strip any path components — prevents directory traversal
	name = filepath.Base(name)
	// Replace all unsafe characters with hyphens
	name = reUnsafe.ReplaceAllString(name, "-")
	// Trim leading/trailing hyphens
	name = strings.Trim(name, "-")
	// Truncate to a reasonable length so the final filename stays manageable
	if len(name) > 64 {
		name = name[:64]
	}
	if name == "" {
		name = "trip"
	}
	return name
}

// uniqueTripFilename returns a filename (without directory) that doesn't
// already exist on disk. It appends a random 6-byte (12-hex-char) suffix to
// ensure uniqueness.
func (s *Server) uniqueTripFilename(base string) (string, error) {
	const maxAttempts = 10
	for i := 0; i < maxAttempts; i++ {
		var buf [6]byte
		if _, err := rand.Read(buf[:]); err != nil {
			return "", fmt.Errorf("failed to generate random suffix: %w", err)
		}
		suffix := hex.EncodeToString(buf[:])
		filename := filepath.Join(s.TripsDir, base+"-"+suffix+".yml")
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return filename, nil
		}
	}
	return "", fmt.Errorf("could not find a unique filename after %d attempts", maxAttempts)
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
