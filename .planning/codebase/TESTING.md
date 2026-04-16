# Testing Patterns

**Analysis Date:** 2026-04-16

## Test Framework

**Runner:**
- Go's built-in `testing` package (no external test runner)
- Run via: `go test ./...` or `go test ./pkg/packinglib`

**Assertion Library:**
- `github.com/stretchr/testify/require` for assertions
- No other assertion libraries; testify is the standard

**Run Commands:**
```bash
go test ./...                    # Run all tests
go test -v ./...                 # Verbose output
go test -run TestName ./...      # Run specific test
go test -cover ./...             # Show coverage
go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out  # HTML coverage report
```

## Test File Organization

**Location:**
- Co-located with source files in same package
- Pattern: `source.go` and `source_test.go` in same directory
- Examples: `pkg/packinglib/item.go` and `pkg/packinglib/item_test.go`

**Naming:**
- Test files: `*_test.go`
- Test functions: `Test<FunctionName>` (e.g., `TestItemAdjustCount`, `TestLookup_Forecast`)
- Sub-tests: `TestName/SubtestName` using `t.Run()`

**Structure:**
```
pkg/packinglib/
├── item.go
├── item_test.go
├── yaml.go
├── yaml_test.go
├── category.go       # no test file
└── context.go        # no test file

pkg/weather/
├── weather.go
└── weather_test.go

cmd/packingweb/
├── main.go
└── e2e_test.go
```

## Test Structure

**Suite Organization:**
From `item_test.go` — table-driven tests with subtests:
```go
func TestItemAdjustCount(t *testing.T) {
	tests := []struct {
		Name          string
		Prerequisites PropertySet
		Mutators      []packMutator
		Context       Context
		Expected      float64
	}{
		{
			Name:          "No adjustments, comes out as 1",
			Prerequisites: PropertySet{},
			Mutators:      []packMutator{},
			Context:       basicContext,
			Expected:      1.0,
		},
		// ... more test cases
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			item := Item{
				prerequisites: tc.Prerequisites,
				mutators:      tc.Mutators,
			}
			item.AdjustCount(&tc.Context)
			require.Equal(t, tc.Expected, item.count)
		})
	}
}
```

**Patterns:**
- Setup: Shared test data (e.g., `basicContext` at line 9) defined at package level
- Teardown: Using `defer` in tests that need cleanup (e.g., `defer server.Close()`)
- Assertion: `require.Equal()`, `require.NoError()`, `require.True()` from testify

**Setup/Teardown from weather_test.go:**
```go
func setupMockServer(t *testing.T, geocodeResults []mockGeocodingResult, dailyMin, dailyMax []float64) *httptest.Server {
	t.Helper()  // Mark as helper function
	
	mux := http.NewServeMux()
	// ... setup handlers
	server := httptest.NewServer(mux)
	
	// Override package-level base URLs
	geocodingBaseURL = server.URL
	forecastBaseURL = server.URL
	archiveBaseURL = server.URL
	
	return server
}

// Called with cleanup:
func TestLookup_Forecast(t *testing.T) {
	server := setupMockServer(t, geoResults, []float64{30.0, 28.5, 32.0}, []float64{55.0, 60.0, 58.0})
	defer server.Close()
	defer restoreBaseURLs()
	
	// ... test code
}
```

## Mocking

**Framework:**
- `net/http/httptest` for HTTP mocking (built-in Go testing library)
- No external mocking framework (Testify is assertion-only)
- Manual mocks for test utilities

**Patterns from weather_test.go:**
```go
// Define mock response structs
type mockGeocodingResult struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Country   string  `json:"country"`
	Admin1    string  `json:"admin1"`
}

// Create httptest.Server with custom handlers
mux := http.NewServeMux()
mux.HandleFunc("/v1/search", func(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		Results []mockGeocodingResult `json:"results"`
	}{Results: geocodeResults}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
})

server := httptest.NewServer(mux)
```

**What to Mock:**
- External HTTP APIs: Open-Meteo geocoding, forecast, archive endpoints
- File I/O when testing creates temp dirs: `os.MkdirTemp()` in e2e_test.go
- Time-dependent functionality: mocked via helper functions (e.g., `localDate()`)

**What NOT to Mock:**
- Business logic internal to the package (Item, Context behavior)
- Registry operations unless testing registry specifically
- Standard library types (only override URL globals in weather tests)

**Example of NOT mocking internal logic (item_test.go):**
```go
// Tests call real Item methods, only mocking the data:
item := Item{
	prerequisites: tc.Prerequisites,
	mutators:      tc.Mutators,  // Real mutator objects, not mocks
}
item.AdjustCount(&tc.Context)   // Calls real method
require.Equal(t, tc.Expected, item.count)
```

## Fixtures and Factories

**Test Data:**
From `item_test.go` — shared context fixture:
```go
var basicContext = Context{
	Name:           "whatever",
	Nights:         3,
	TemperatureMin: 50,
	TemperatureMax: 80,
	Properties:     PropertySet{"prop1": true, "prop2": true},
}

// Reused across multiple tests
func TestItemAdjustCount(t *testing.T) {
	// ... uses basicContext
}

func TestItemString(t *testing.T) {
	// ... also uses basicContext
}
```

From `yaml_test.go` — factory function:
```go
func populateRegistry(r Registry) {
	r.RegisterProperty("Business", "business cat")
	r.RegisterProperty("Flight", "wheee")
	// ...
	r.RegisterItems("business", []*Item{
		NewItem("work laptop", []string{"Business"}, nil),
		NewItem("socks", []string{"Flight"}, []string{"camping"}).Consumable(1),
	})
}

// Used in tests:
var r Registry = NewStructRegistry()
populateRegistry(r)
got, err := yt.AsTrip(r)
```

From `weather_test.go` — helper factory:
```go
func setupMockServer(t *testing.T, geocodeResults []mockGeocodingResult, 
                     dailyMin, dailyMax []float64) *httptest.Server {
	t.Helper()
	// ... returns configured server
}

// Helper to restore state:
func restoreBaseURLs() {
	geocodingBaseURL = "https://geocoding-api.open-meteo.com"
	forecastBaseURL = "https://api.open-meteo.com"
	archiveBaseURL = "https://archive-api.open-meteo.com"
}
```

**Location:**
- Fixtures defined at package level in `_test.go` files
- Factories defined as functions in `_test.go` files (prefixed `setup`, `populate`, `restore`)
- No separate fixtures directory; all in test files

## Coverage

**Requirements:**
- No `.coveragerc` or CI enforcer found
- Coverage appears encouraged but not strictly required
- Some coverage gaps exist (e.g., category.go, context.go have no tests)

**View Coverage:**
```bash
go test -cover ./pkg/packinglib
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**Gaps identified:**
- `pkg/packinglib/category.go` - Not tested
- `pkg/packinglib/context.go` - Not tested (except indirectly via yaml_test.go)
- `pkg/packinglib/property.go` - Not tested
- `pkg/server/persist.go` - Not tested (file I/O operations)
- `pkg/items/*.go` - Not tested (data registration files)

## Test Types

**Unit Tests:**
- Scope: Individual functions and methods
- Approach: Table-driven tests with `t.Run()` subtests
- Examples: `TestItemAdjustCount`, `TestItemString`, `TestItemConstruction`, `TestYamlItemList`
- Coverage: Business logic for Item, Context, Yaml serialization, weather API calls

**Integration Tests:**
- Scope: Multiple components working together
- Approach: Setup test data, call multiple functions, verify state
- Example: `TestYamlTrip` (yaml_test.go line 95) — creates a Trip from YAML, checks it matches expected state, serializes back to YAML
- Coverage: Item + Mutator combinations, Trip creation + properties, Registry + Trip interaction

**E2E Tests:**
- Framework: `testing.T` with `os/exec` to launch actual binary
- Approach: Start server, make HTTP requests, verify responses
- Location: `cmd/packingweb/e2e_test.go`
- Coverage: HTTP API contracts, server startup, persistence
- Example: `TestE2E` (line 18) — builds binary, starts server on random port, runs suite of API tests

**E2E Test Pattern:**
```go
func TestE2E(t *testing.T) {
	// Create temp environment
	tmpDir, err := os.MkdirTemp("", "packingweb-test-*")
	defer os.RemoveAll(tmpDir)
	
	// Find unused port
	port, err := getFreePort()
	
	// Build binary
	binaryPath := filepath.Join(tmpDir, "packingweb")
	cmd := exec.Command("go", "build", "-o", binaryPath)
	
	// Start server
	serverCmd := exec.Command(binaryPath, "-trips", tripsDir, "-static", staticDir, "-port", fmt.Sprintf("%d", port))
	serverCmd.Start()
	defer serverCmd.Process.Signal(os.Interrupt)  // Graceful shutdown
	
	// Wait for readiness
	baseURL := fmt.Sprintf("http://localhost:%d", port)
	if !waitForServer(baseURL, 5*time.Second) { t.Fatal("timeout") }
	
	// Run sub-tests
	t.Run("ListTripsEmpty", func(t *testing.T) { ... })
	t.Run("CreateTrip", func(t *testing.T) { ... })
}
```

## Common Patterns

**Async Testing:**
From `e2e_test.go` — waiting for server readiness:
```go
// Wait for server startup
if !waitForServer(baseURL, 5*time.Second) {
	t.Fatal("Server failed to start within timeout")
}

// With timeout and cleanup:
done := make(chan struct{})
go func() {
	_ = serverCmd.Wait()
	close(done)
}()

select {
case <-done:
	// Graceful shutdown completed
case <-time.After(5 * time.Second):
	// Fallback to kill
	_ = serverCmd.Process.Kill()
}
```

**Error Testing:**
From `weather_test.go` — testing error conditions:
```go
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

func TestLookup_WeatherAPIError(t *testing.T) {
	// Setup mock that returns error
	mux.HandleFunc("/v1/forecast", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	
	_, err := Lookup("Denver", start, end)
	if err == nil {
		t.Fatal("expected error when weather API fails, got nil")
	}
	if !strings.Contains(err.Error(), "weather fetch failed") {
		t.Errorf("expected 'weather fetch failed' in error, got: %v", err)
	}
}
```

**Testing Slice/Map Results:**
From `weather_test.go`:
```go
func TestGeocodeSuggestions_Basic(t *testing.T) {
	geoResults := []mockGeocodingResult{
		{Name: "Portland", Latitude: 45.52, Longitude: -122.68, Country: "United States", Admin1: "Oregon"},
		{Name: "Portland", Latitude: 43.66, Longitude: -70.26, Country: "United States", Admin1: "Maine"},
	}
	
	suggestions, err := GeocodeSuggestions("Portland")
	if err != nil { t.Fatalf(...) }
	if len(suggestions) != 2 { t.Fatalf(...) }
	if suggestions[0].Display != "Portland, Oregon, United States" {
		t.Errorf("unexpected display: %q", suggestions[0].Display)
	}
}
```

## Test Documentation

**Comments:**
- Helper functions marked with `t.Helper()` to improve error reporting
- Complex test setup documented with inline comments
- Table test case names describe the scenario: `"No adjustments, comes out as 1"`, `"Temperator Mutator min denies"`

**Example from weather_test.go:**
```go
// setupMockServer creates an httptest.Server that handles geocoding, forecast,
// and archive endpoints. It returns the server and overrides the package-level
// base URLs to point to it. Callers are responsible for restoring the original
// base URLs after the test completes.
func setupMockServer(t *testing.T, geocodeResults []mockGeocodingResult, 
                     dailyMin, dailyMax []float64) *httptest.Server {
	t.Helper()
	// ...
}
```

---

*Testing analysis: 2026-04-16*
