package weather

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// mockGeocodingResult is used to build mock geocoding API responses.
type mockGeocodingResult struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Country   string  `json:"country"`
	Admin1    string  `json:"admin1"`
}

// mockDailyWeather is used to build mock daily weather API responses.
type mockDailyWeather struct {
	Daily struct {
		TemperatureMax []float64 `json:"temperature_2m_max"`
		TemperatureMin []float64 `json:"temperature_2m_min"`
	} `json:"daily"`
}

// setupMockServer creates an httptest.Server that handles geocoding, forecast,
// and archive endpoints. It returns the server and a cleanup function that
// restores the original base URLs.
func setupMockServer(t *testing.T, geocodeResults []mockGeocodingResult, dailyMin, dailyMax []float64) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()

	// Geocoding endpoint
	mux.HandleFunc("/v1/search", func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Results []mockGeocodingResult `json:"results"`
		}{Results: geocodeResults}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	// Forecast endpoint
	mux.HandleFunc("/v1/forecast", func(w http.ResponseWriter, r *http.Request) {
		var resp mockDailyWeather
		resp.Daily.TemperatureMin = dailyMin
		resp.Daily.TemperatureMax = dailyMax
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	// Historical archive endpoint
	mux.HandleFunc("/v1/archive", func(w http.ResponseWriter, r *http.Request) {
		var resp mockDailyWeather
		resp.Daily.TemperatureMin = dailyMin
		resp.Daily.TemperatureMax = dailyMax
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	server := httptest.NewServer(mux)

	// Override package-level base URLs to point at the mock server.
	geocodingBaseURL = server.URL
	forecastBaseURL = server.URL
	archiveBaseURL = server.URL

	return server
}

// restoreBaseURLs resets the base URLs to their production values.
func restoreBaseURLs() {
	geocodingBaseURL = "https://geocoding-api.open-meteo.com"
	forecastBaseURL = "https://api.open-meteo.com"
	archiveBaseURL = "https://archive-api.open-meteo.com"
}

// --- Lookup tests ---

func TestLookup_Forecast(t *testing.T) {
	geoResults := []mockGeocodingResult{
		{Name: "Denver", Latitude: 39.7392, Longitude: -104.9903, Country: "United States", Admin1: "Colorado"},
	}
	server := setupMockServer(t, geoResults, []float64{30.0, 28.5, 32.0}, []float64{55.0, 60.0, 58.0})
	defer server.Close()
	defer restoreBaseURLs()

	// Dates within the 16-day forecast window.
	start := localDate(time.Now()).AddDate(0, 0, 1)
	end := start.AddDate(0, 0, 3)

	result, err := Lookup("Denver", start, end)
	if err != nil {
		t.Fatalf("Lookup returned error: %v", err)
	}

	if result.Source != "forecast" {
		t.Errorf("expected source 'forecast', got %q", result.Source)
	}
	if result.Location != "Denver, Colorado, United States" {
		t.Errorf("unexpected location: %q", result.Location)
	}
	if result.Nights != 3 {
		t.Errorf("expected 3 nights, got %d", result.Nights)
	}
	if result.TemperatureMin != 29 { // round(28.5)
		t.Errorf("expected min temp 29, got %d", result.TemperatureMin)
	}
	if result.TemperatureMax != 60 {
		t.Errorf("expected max temp 60, got %d", result.TemperatureMax)
	}
}

func TestLookup_RejectsPastDates(t *testing.T) {
	// No mock server needed — Lookup should reject before making any HTTP call.
	start := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 7, 5, 0, 0, 0, 0, time.UTC)

	_, err := Lookup("Berlin", start, end)
	if err == nil {
		t.Fatal("expected error for past dates, got nil")
	}
	if !strings.Contains(err.Error(), "start date must be today or in the future") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestLookup_Typical(t *testing.T) {
	geoResults := []mockGeocodingResult{
		{Name: "Tokyo", Latitude: 35.6762, Longitude: 139.6503, Country: "Japan", Admin1: "Tokyo"},
	}
	server := setupMockServer(t, geoResults, []float64{60.0, 62.0}, []float64{85.0, 88.0})
	defer server.Close()
	defer restoreBaseURLs()

	// Dates far in the future (>16 days) → typical (previous year).
	start := localDate(time.Now()).AddDate(0, 3, 0)
	end := start.AddDate(0, 0, 5)

	result, err := Lookup("Tokyo", start, end)
	if err != nil {
		t.Fatalf("Lookup returned error: %v", err)
	}

	if result.Source != "typical" {
		t.Errorf("expected source 'typical', got %q", result.Source)
	}
	if result.Nights != 5 {
		t.Errorf("expected 5 nights, got %d", result.Nights)
	}
	if result.TemperatureMin != 60 {
		t.Errorf("expected min temp 60, got %d", result.TemperatureMin)
	}
	if result.TemperatureMax != 88 {
		t.Errorf("expected max temp 88, got %d", result.TemperatureMax)
	}
}

func TestLookupByCoords_Forecast(t *testing.T) {
	// LookupByCoords doesn't need geocoding results — it skips geocoding.
	server := setupMockServer(t, nil, []float64{30.0, 28.5, 32.0}, []float64{55.0, 60.0, 58.0})
	defer server.Close()
	defer restoreBaseURLs()

	start := localDate(time.Now()).AddDate(0, 0, 1)
	end := start.AddDate(0, 0, 3)

	result, err := LookupByCoords(42.36, -73.59, "Chatham, New York, United States", start, end)
	if err != nil {
		t.Fatalf("LookupByCoords returned error: %v", err)
	}
	if result.Source != "forecast" {
		t.Errorf("expected source 'forecast', got %q", result.Source)
	}
	if result.Location != "Chatham, New York, United States" {
		t.Errorf("unexpected location: %q", result.Location)
	}
	if result.Nights != 3 {
		t.Errorf("expected 3 nights, got %d", result.Nights)
	}
	if result.TemperatureMin != 29 {
		t.Errorf("expected min temp 29, got %d", result.TemperatureMin)
	}
	if result.TemperatureMax != 60 {
		t.Errorf("expected max temp 60, got %d", result.TemperatureMax)
	}
}

func TestLookupByCoords_RejectsPastDates(t *testing.T) {
	start := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 7, 5, 0, 0, 0, 0, time.UTC)

	_, err := LookupByCoords(42.36, -73.59, "Chatham, New York, United States", start, end)
	if err == nil {
		t.Fatal("expected error for past dates, got nil")
	}
}

func TestLookup_GeocodeFail(t *testing.T) {
	// Empty geocoding results → error.
	server := setupMockServer(t, nil, nil, nil)
	defer server.Close()
	defer restoreBaseURLs()

	start := localDate(time.Now()).AddDate(0, 0, 1)
	end := start.AddDate(0, 0, 2)

	_, err := Lookup("Nonexistent Place XYZ", start, end)
	if err == nil {
		t.Fatal("expected error for unknown location, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' in error, got: %v", err)
	}
}

func TestLookup_MinNightsClamp(t *testing.T) {
	geoResults := []mockGeocodingResult{
		{Name: "Boston", Latitude: 42.36, Longitude: -71.06, Country: "United States", Admin1: "Massachusetts"},
	}
	server := setupMockServer(t, geoResults, []float64{40.0}, []float64{65.0})
	defer server.Close()
	defer restoreBaseURLs()

	// Same start and end date → nights should be clamped to 1.
	d := localDate(time.Now()).AddDate(0, 0, 1)
	result, err := Lookup("Boston", d, d)
	if err != nil {
		t.Fatalf("Lookup returned error: %v", err)
	}
	if result.Nights != 1 {
		t.Errorf("expected nights clamped to 1, got %d", result.Nights)
	}
}

func TestLookup_RejectsYesterdayLocalDate(t *testing.T) {
	start := localDate(time.Now()).AddDate(0, 0, -1)
	end := start.AddDate(0, 0, 2)

	_, err := Lookup("Berlin", start, end)
	if err == nil {
		t.Fatal("expected error for yesterday's local date, got nil")
	}
	if !strings.Contains(err.Error(), "start date must be today or in the future") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// --- GeocodeSuggestions tests ---

func TestGeocodeSuggestions_Basic(t *testing.T) {
	geoResults := []mockGeocodingResult{
		{Name: "Portland", Latitude: 45.52, Longitude: -122.68, Country: "United States", Admin1: "Oregon"},
		{Name: "Portland", Latitude: 43.66, Longitude: -70.26, Country: "United States", Admin1: "Maine"},
	}
	server := setupMockServer(t, geoResults, nil, nil)
	defer server.Close()
	defer restoreBaseURLs()

	suggestions, err := GeocodeSuggestions("Portland")
	if err != nil {
		t.Fatalf("GeocodeSuggestions returned error: %v", err)
	}
	if len(suggestions) != 2 {
		t.Fatalf("expected 2 suggestions, got %d", len(suggestions))
	}
	if suggestions[0].Display != "Portland, Oregon, United States" {
		t.Errorf("unexpected display: %q", suggestions[0].Display)
	}
	if suggestions[1].Display != "Portland, Maine, United States" {
		t.Errorf("unexpected display: %q", suggestions[1].Display)
	}
}

func TestGeocodeSuggestions_CommaFilter(t *testing.T) {
	geoResults := []mockGeocodingResult{
		{Name: "Chatham", Latitude: 36.83, Longitude: -79.40, Country: "United States", Admin1: "Virginia"},
		{Name: "Chatham", Latitude: 51.38, Longitude: 0.52, Country: "United Kingdom", Admin1: "England"},
		{Name: "Chatham", Latitude: 42.36, Longitude: -73.59, Country: "United States", Admin1: "New York"},
		{Name: "Chatham", Latitude: 42.41, Longitude: -82.19, Country: "Canada", Admin1: "Ontario"},
	}
	server := setupMockServer(t, geoResults, nil, nil)
	defer server.Close()
	defer restoreBaseURLs()

	suggestions, err := GeocodeSuggestions("Chatham, Virginia")
	if err != nil {
		t.Fatalf("GeocodeSuggestions returned error: %v", err)
	}
	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
	}
	if suggestions[0].Admin1 != "Virginia" {
		t.Errorf("expected Virginia, got %q", suggestions[0].Admin1)
	}
}

func TestGeocodeSuggestions_AbbreviationFilter(t *testing.T) {
	geoResults := []mockGeocodingResult{
		{Name: "Chatham", Latitude: 36.83, Longitude: -79.40, Country: "United States", Admin1: "Virginia"},
		{Name: "Chatham", Latitude: 51.38, Longitude: 0.52, Country: "United Kingdom", Admin1: "England"},
		{Name: "Chatham", Latitude: 42.36, Longitude: -73.59, Country: "United States", Admin1: "New York"},
		{Name: "Chatham", Latitude: 42.41, Longitude: -82.19, Country: "Canada", Admin1: "Ontario"},
	}
	server := setupMockServer(t, geoResults, nil, nil)
	defer server.Close()
	defer restoreBaseURLs()

	suggestions, err := GeocodeSuggestions("Chatham, ny")
	if err != nil {
		t.Fatalf("GeocodeSuggestions returned error: %v", err)
	}
	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
	}
	if suggestions[0].Admin1 != "New York" {
		t.Errorf("expected New York, got %q", suggestions[0].Admin1)
	}
}

func TestGeocodeSuggestions_CanadianAbbreviation(t *testing.T) {
	geoResults := []mockGeocodingResult{
		{Name: "Vancouver", Latitude: 49.25, Longitude: -123.12, Country: "Canada", Admin1: "British Columbia"},
		{Name: "Vancouver", Latitude: 45.63, Longitude: -122.66, Country: "United States", Admin1: "Washington"},
	}
	server := setupMockServer(t, geoResults, nil, nil)
	defer server.Close()
	defer restoreBaseURLs()

	suggestions, err := GeocodeSuggestions("Vancouver, bc")
	if err != nil {
		t.Fatalf("GeocodeSuggestions returned error: %v", err)
	}
	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
	}
	if suggestions[0].Admin1 != "British Columbia" {
		t.Errorf("expected British Columbia, got %q", suggestions[0].Admin1)
	}
}

func TestGeocodeSuggestions_EmptyQuery(t *testing.T) {
	server := setupMockServer(t, nil, nil, nil)
	defer server.Close()
	defer restoreBaseURLs()

	suggestions, err := GeocodeSuggestions("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if suggestions != nil {
		t.Errorf("expected nil suggestions for empty query, got %v", suggestions)
	}
}

func TestGeocodeSuggestions_MaxFiveResults(t *testing.T) {
	geoResults := make([]mockGeocodingResult, 8)
	for i := range geoResults {
		geoResults[i] = mockGeocodingResult{
			Name:      "Springfield",
			Latitude:  39.0 + float64(i),
			Longitude: -89.0,
			Country:   "United States",
			Admin1:    "State" + string(rune('A'+i)),
		}
	}
	server := setupMockServer(t, geoResults, nil, nil)
	defer server.Close()
	defer restoreBaseURLs()

	suggestions, err := GeocodeSuggestions("Springfield")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(suggestions) != 5 {
		t.Errorf("expected max 5 suggestions, got %d", len(suggestions))
	}
}

func TestGeocodeSuggestions_NoAdmin1(t *testing.T) {
	geoResults := []mockGeocodingResult{
		{Name: "Singapore", Latitude: 1.35, Longitude: 103.82, Country: "Singapore"},
	}
	server := setupMockServer(t, geoResults, nil, nil)
	defer server.Close()
	defer restoreBaseURLs()

	suggestions, err := GeocodeSuggestions("Singapore")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
	}
	if suggestions[0].Display != "Singapore, Singapore" {
		t.Errorf("unexpected display: %q", suggestions[0].Display)
	}
}

// --- matchesAdmin1Filter tests ---

func TestMatchesAdmin1Filter_PrefixMatch(t *testing.T) {
	tests := []struct {
		admin1, country, filter string
		want                    bool
	}{
		{"New York", "United States", "new", true},
		{"California", "United States", "cal", true},
		{"Ontario", "Canada", "on", true},
		{"British Columbia", "Canada", "british", true},
		{"Texas", "United States", "florida", false},
		{"", "United Kingdom", "united", true},
	}

	for _, tt := range tests {
		got := matchesAdmin1Filter(tt.admin1, tt.country, tt.filter)
		if got != tt.want {
			t.Errorf("matchesAdmin1Filter(%q, %q, %q) = %v, want %v",
				tt.admin1, tt.country, tt.filter, got, tt.want)
		}
	}
}

func TestMatchesAdmin1Filter_Abbreviation(t *testing.T) {
	tests := []struct {
		admin1, country, filter string
		want                    bool
	}{
		{"New York", "United States", "ny", true},
		{"California", "United States", "ca", true},
		{"British Columbia", "Canada", "bc", true},
		{"Oregon", "United States", "or", true},
		{"Texas", "United States", "ny", false},
		{"Nova Scotia", "Canada", "ns", true},
		{"Quebec", "Canada", "qc", true},
	}

	for _, tt := range tests {
		got := matchesAdmin1Filter(tt.admin1, tt.country, tt.filter)
		if got != tt.want {
			t.Errorf("matchesAdmin1Filter(%q, %q, %q) = %v, want %v",
				tt.admin1, tt.country, tt.filter, got, tt.want)
		}
	}
}

// --- Error handling tests ---

func TestLookup_WeatherAPIError(t *testing.T) {
	geoResults := []mockGeocodingResult{
		{Name: "Denver", Latitude: 39.74, Longitude: -104.99, Country: "United States", Admin1: "Colorado"},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/search", func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Results []mockGeocodingResult `json:"results"`
		}{Results: geoResults}
		json.NewEncoder(w).Encode(resp)
	})
	mux.HandleFunc("/v1/forecast", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	mux.HandleFunc("/v1/archive", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	geocodingBaseURL = server.URL
	forecastBaseURL = server.URL
	archiveBaseURL = server.URL
	defer restoreBaseURLs()

	start := time.Now().UTC().Truncate(24*time.Hour).AddDate(0, 0, 1)
	end := start.AddDate(0, 0, 2)

	_, err := Lookup("Denver", start, end)
	if err == nil {
		t.Fatal("expected error when weather API fails, got nil")
	}
	if !strings.Contains(err.Error(), "weather fetch failed") {
		t.Errorf("expected 'weather fetch failed' in error, got: %v", err)
	}
}

func TestLookup_EmptyWeatherData(t *testing.T) {
	geoResults := []mockGeocodingResult{
		{Name: "Denver", Latitude: 39.74, Longitude: -104.99, Country: "United States", Admin1: "Colorado"},
	}
	// Empty temperature arrays → should error.
	server := setupMockServer(t, geoResults, []float64{}, []float64{})
	defer server.Close()
	defer restoreBaseURLs()

	start := time.Now().UTC().Truncate(24*time.Hour).AddDate(0, 0, 1)
	end := start.AddDate(0, 0, 2)

	_, err := Lookup("Denver", start, end)
	if err == nil {
		t.Fatal("expected error for empty weather data, got nil")
	}
}

func TestGeocodeSuggestions_ServerError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/search", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad gateway", http.StatusBadGateway)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	geocodingBaseURL = server.URL
	defer restoreBaseURLs()

	_, err := GeocodeSuggestions("Denver")
	if err == nil {
		t.Fatal("expected error when geocoding API returns error, got nil")
	}
}
