package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestE2E runs end-to-end tests against the actual packingweb binary
func TestE2E(t *testing.T) {
	// Create temporary directory for trips
	tmpDir, err := os.MkdirTemp("", "packingweb-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tripsDir := filepath.Join(tmpDir, "trips")
	if err := os.MkdirAll(tripsDir, 0755); err != nil {
		t.Fatalf("Failed to create trips directory: %v", err)
	}

	// Find an unused port
	port, err := getFreePort()
	if err != nil {
		t.Fatalf("Failed to find free port: %v", err)
	}

	// Build the binary
	binaryPath := filepath.Join(tmpDir, "packingweb")
	cmd := exec.Command("go", "build", "-o", binaryPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\nOutput: %s", err, output)
	}

	// Determine static directory (relative to test location)
	staticDir := "../../static"
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		t.Fatalf("Static directory not found at %s", staticDir)
	}

	// Start the server
	serverCmd := exec.Command(binaryPath)
	serverCmd.Env = append(os.Environ(),
		fmt.Sprintf("PORT=%d", port),
		fmt.Sprintf("TRIPS_DIR=%s", tripsDir),
		fmt.Sprintf("STATIC_DIR=%s", staticDir),
	)

	// Capture output for debugging
	serverCmd.Stdout = os.Stdout
	serverCmd.Stderr = os.Stderr

	if err := serverCmd.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer serverCmd.Process.Kill()

	// Wait for server to be ready
	baseURL := fmt.Sprintf("http://localhost:%d", port)
	if !waitForServer(baseURL, 5*time.Second) {
		t.Fatal("Server failed to start within timeout")
	}

	t.Logf("Server started on port %d", port)

	// Run tests
	t.Run("ListTripsEmpty", func(t *testing.T) {
		testListTripsEmpty(t, baseURL)
	})

	t.Run("CreateTrip", func(t *testing.T) {
		testCreateTrip(t, baseURL, tripsDir)
	})

	t.Run("GetTrip", func(t *testing.T) {
		testGetTrip(t, baseURL)
	})

	t.Run("GetProperties", func(t *testing.T) {
		testGetProperties(t, baseURL)
	})

	t.Run("GetTripItems", func(t *testing.T) {
		testGetTripItems(t, baseURL)
	})

	t.Run("ToggleItem", func(t *testing.T) {
		testToggleItem(t, baseURL)
	})

	t.Run("UpdateTrip", func(t *testing.T) {
		testUpdateTrip(t, baseURL)
	})
}

// getFreePort finds an available port on the system
func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// waitForServer polls the server until it's ready or timeout
func waitForServer(baseURL string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(baseURL + "/api/trips")
		if err == nil {
			resp.Body.Close()
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}
	return false
}

// testListTripsEmpty verifies empty trip list
func testListTripsEmpty(t *testing.T, baseURL string) {
	resp, err := http.Get(baseURL + "/api/trips")
	if err != nil {
		t.Fatalf("Failed to get trips: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	trips, ok := result["trips"].([]interface{})
	if !ok {
		t.Fatal("Expected trips array in response")
	}

	if len(trips) != 0 {
		t.Errorf("Expected 0 trips, got %d", len(trips))
	}
}

// testCreateTrip creates a new trip
func testCreateTrip(t *testing.T, baseURL string, tripsDir string) {
	payload := strings.NewReader(`{
		"name": "test-trip",
		"nights": 3,
		"temperatureMin": 60,
		"temperatureMax": 80,
		"properties": ["Hiking", "Camping"]
	}`)

	resp, err := http.Post(baseURL+"/api/trips", "application/json", payload)
	if err != nil {
		t.Fatalf("Failed to create trip: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200, got %d. Body: %s", resp.StatusCode, body)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if success, ok := result["success"].(bool); !ok || !success {
		t.Error("Expected success: true")
	}

	// Verify file was created
	tripFile := filepath.Join(tripsDir, "test-trip.yml")
	if _, err := os.Stat(tripFile); os.IsNotExist(err) {
		t.Errorf("Trip file was not created at %s", tripFile)
	}
}

// testGetTrip retrieves a trip
func testGetTrip(t *testing.T, baseURL string) {
	resp, err := http.Get(baseURL + "/api/trips/test-trip")
	if err != nil {
		t.Fatalf("Failed to get trip: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200, got %d. Body: %s", resp.StatusCode, body)
	}

	var trip map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&trip); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if name := trip["name"]; name != "test-trip" {
		t.Errorf("Expected name 'test-trip', got %v", name)
	}

	if nights := trip["nights"]; nights != float64(3) {
		t.Errorf("Expected 3 nights, got %v", nights)
	}

	if tempMin := trip["temperatureMin"]; tempMin != float64(60) {
		t.Errorf("Expected temperatureMin 60, got %v", tempMin)
	}
}

// testGetProperties retrieves all properties
func testGetProperties(t *testing.T, baseURL string) {
	resp, err := http.Get(baseURL + "/api/properties")
	if err != nil {
		t.Fatalf("Failed to get properties: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	properties, ok := result["properties"].([]interface{})
	if !ok {
		t.Fatal("Expected properties array in response")
	}

	if len(properties) == 0 {
		t.Error("Expected at least one property")
	}
}

// testGetTripItems retrieves trip items
func testGetTripItems(t *testing.T, baseURL string) {
	resp, err := http.Get(baseURL + "/api/trips/test-trip/items")
	if err != nil {
		t.Fatalf("Failed to get items: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200, got %d. Body: %s", resp.StatusCode, body)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	categories, ok := result["categories"].([]interface{})
	if !ok {
		t.Fatal("Expected categories array in response")
	}

	if len(categories) == 0 {
		t.Error("Expected at least one category")
	}
}

// testToggleItem toggles an item's packed status
func testToggleItem(t *testing.T, baseURL string) {
	// First, get an item code
	resp, err := http.Get(baseURL + "/api/trips/test-trip/items")
	if err != nil {
		t.Fatalf("Failed to get items: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	categories := result["categories"].([]interface{})
	if len(categories) == 0 {
		t.Skip("No categories available to test")
	}

	category := categories[0].(map[string]interface{})
	items := category["items"].([]interface{})
	if len(items) == 0 {
		t.Skip("No items available to test")
	}

	item := items[0].(map[string]interface{})
	code := item["code"].(string)
	initialPacked := item["packed"].(bool)

	// Toggle the item
	toggleURL := fmt.Sprintf("%s/api/trips/test-trip/items/%s/toggle", baseURL, code)
	toggleResp, err := http.Post(toggleURL, "application/json", nil)
	if err != nil {
		t.Fatalf("Failed to toggle item: %v", err)
	}
	defer toggleResp.Body.Close()

	if toggleResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(toggleResp.Body)
		t.Fatalf("Expected status 200, got %d. Body: %s", toggleResp.StatusCode, body)
	}

	// Verify the toggle worked
	resp2, err := http.Get(baseURL + "/api/trips/test-trip/items")
	if err != nil {
		t.Fatalf("Failed to get items after toggle: %v", err)
	}
	defer resp2.Body.Close()

	var result2 map[string]interface{}
	if err := json.NewDecoder(resp2.Body).Decode(&result2); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	categories2 := result2["categories"].([]interface{})
	category2 := categories2[0].(map[string]interface{})
	items2 := category2["items"].([]interface{})
	item2 := items2[0].(map[string]interface{})
	newPacked := item2["packed"].(bool)

	if newPacked == initialPacked {
		t.Errorf("Item packed status did not toggle (stayed %v)", initialPacked)
	}
}

// testUpdateTrip updates trip settings
func testUpdateTrip(t *testing.T, baseURL string) {
	payload := strings.NewReader(`{
		"nights": 5,
		"temperatureMin": 50,
		"temperatureMax": 90
	}`)

	req, err := http.NewRequest("PUT", baseURL+"/api/trips/test-trip/update", payload)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to update trip: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200, got %d. Body: %s", resp.StatusCode, body)
	}

	// Verify the update
	getResp, err := http.Get(baseURL + "/api/trips/test-trip")
	if err != nil {
		t.Fatalf("Failed to get trip after update: %v", err)
	}
	defer getResp.Body.Close()

	var trip map[string]interface{}
	if err := json.NewDecoder(getResp.Body).Decode(&trip); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if nights := trip["nights"]; nights != float64(5) {
		t.Errorf("Expected 5 nights after update, got %v", nights)
	}

	if tempMin := trip["temperatureMin"]; tempMin != float64(50) {
		t.Errorf("Expected temperatureMin 50 after update, got %v", tempMin)
	}

	if tempMax := trip["temperatureMax"]; tempMax != float64(90) {
		t.Errorf("Expected temperatureMax 90 after update, got %v", tempMax)
	}
}
