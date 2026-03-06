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

	// Start the server with command-line flags
	serverCmd := exec.Command(binaryPath,
		"-trips", tripsDir,
		"-static", staticDir,
		"-port", fmt.Sprintf("%d", port),
	)

	// Capture output for debugging
	serverCmd.Stdout = os.Stdout
	serverCmd.Stderr = os.Stderr

	if err := serverCmd.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer func() {
		if serverCmd.Process == nil {
			return
		}
		// Try to shut down the server gracefully to allow cleanup of background goroutines.
		_ = serverCmd.Process.Signal(os.Interrupt)

		done := make(chan struct{})
		go func() {
			// Wait for the process to exit; ignore the error here since we're in cleanup.
			_ = serverCmd.Wait()
			close(done)
		}()

		select {
		case <-done:
			// Graceful shutdown completed.
		case <-time.After(5 * time.Second):
			// Fallback to forceful kill if the process does not exit in time.
			_ = serverCmd.Process.Kill()
		}
	}()
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

	var tripKey string
	t.Run("CreateTrip", func(t *testing.T) {
		tripKey = testCreateTrip(t, baseURL, tripsDir)
	})

	t.Run("GetTrip", func(t *testing.T) {
		testGetTrip(t, baseURL, tripKey)
	})

	t.Run("GetProperties", func(t *testing.T) {
		testGetProperties(t, baseURL)
	})

	t.Run("GetTripItems", func(t *testing.T) {
		testGetTripItems(t, baseURL, tripKey)
	})

	t.Run("ToggleItem", func(t *testing.T) {
		testToggleItem(t, baseURL, tripKey)
	})

	t.Run("GetTripProperties", func(t *testing.T) {
		testGetTripProperties(t, baseURL, tripKey)
	})

	t.Run("GetAllProperties", func(t *testing.T) {
		testGetAllProperties(t, baseURL)
	})

	t.Run("GetTripItems", func(t *testing.T) {
		testGetTripItems(t, baseURL, tripKey)
	})

	t.Run("GetTripProperties", func(t *testing.T) {
		testGetTripProperties(t, baseURL, tripKey)
	})

	t.Run("ToggleItem", func(t *testing.T) {
		testToggleItem(t, baseURL, tripKey)
	})

	t.Run("ToggleTripProperty", func(t *testing.T) {
		testToggleTripProperty(t, baseURL, tripKey)
	})

	t.Run("UpdateTrip", func(t *testing.T) {
		testUpdateTrip(t, baseURL, tripKey)
	})

	t.Run("RenameTrip", func(t *testing.T) {
		testRenameTrip(t, baseURL, tripKey)
	})
}

func testGetAllProperties(t *testing.T, baseURL string) {
	resp, err := http.Get(baseURL + "/api/properties")
	if err != nil {
		t.Fatalf("failed to GET all properties: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /api/properties returned status %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var body interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil && err != io.EOF {
		t.Fatalf("failed to decode /api/properties response: %v", err)
	}
}

func testGetTripProperties(t *testing.T, baseURL, tripKey string) {
	url := fmt.Sprintf("%s/api/trips/%s/properties", baseURL, tripKey)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("failed to GET trip properties: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET %%s returned status %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var body interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil && err != io.EOF {
		t.Fatalf("failed to decode trip properties response: %v", err)
	}
}

func testToggleTripProperty(t *testing.T, baseURL, tripKey string) {
	url := fmt.Sprintf("%s/api/trips/%s/properties/toggle", baseURL, tripKey)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		t.Fatalf("failed to create toggle trip property request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to POST trip property toggle: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		t.Fatalf("POST %%s returned status %d, want %d or %d", resp.StatusCode, http.StatusOK, http.StatusNoContent)
	}
}

func testRenameTrip(t *testing.T, baseURL, tripKey string) {
	newName := "renamed-trip"
	url := fmt.Sprintf("%s/api/trips/%s/update", baseURL, tripKey)
	body := fmt.Sprintf(`{"name":%q}`, newName)

	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(body))
	if err != nil {
		t.Fatalf("failed to create rename trip request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to PUT trip rename: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.Fatalf("PUT %%s returned status %d, want 2xx", resp.StatusCode)
	}

	// Verify that the trip's name was updated by fetching it again.
	getURL := fmt.Sprintf("%s/api/trips/%s", baseURL, tripKey)
	getResp, err := http.Get(getURL)
	if err != nil {
		t.Fatalf("failed to GET trip after rename: %v", err)
	}
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("GET %%s after rename returned status %d, want %d", getResp.StatusCode, http.StatusOK)
	}

	var trip map[string]interface{}
	if err := json.NewDecoder(getResp.Body).Decode(&trip); err != nil {
		t.Fatalf("failed to decode trip after rename: %v", err)
	}

	if name, ok := trip["name"].(string); !ok || name != newName {
		t.Fatalf("trip name after rename = %v, want %q", trip["name"], newName)
	}
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

// testCreateTrip creates a new trip and returns the trip name
func testCreateTrip(t *testing.T, baseURL string, tripsDir string) string {
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

	// The API returns the base filename (e.g. "test-trip-abc123.yml")
	filename, _ := result["filename"].(string)
	if filename == "" {
		t.Fatal("Expected filename in response")
	}

	// Verify file was created on disk
	tripFile := filepath.Join(tripsDir, filename)
	if _, err := os.Stat(tripFile); os.IsNotExist(err) {
		t.Errorf("Trip file was not created at %s", tripFile)
	}

	// Return the trip name (used in API URLs, not the filename)
	return "test-trip"
}

// testGetTrip retrieves a trip
func testGetTrip(t *testing.T, baseURL string, tripKey string) {
	resp, err := http.Get(baseURL + "/api/trips/" + tripKey)
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
func testGetTripItems(t *testing.T, baseURL string, tripKey string) {
	resp, err := http.Get(baseURL + "/api/trips/" + tripKey + "/items")
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
func testToggleItem(t *testing.T, baseURL string, tripKey string) {
	// First, get an item code
	resp, err := http.Get(baseURL + "/api/trips/" + tripKey + "/items")
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
	toggleURL := fmt.Sprintf("%s/api/trips/%s/items/%s/toggle", baseURL, tripKey, code)
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
	resp2, err := http.Get(baseURL + "/api/trips/" + tripKey + "/items")
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
func testUpdateTrip(t *testing.T, baseURL string, tripKey string) {
	payload := strings.NewReader(`{
		"nights": 5,
		"temperatureMin": 50,
		"temperatureMax": 90
	}`)

	req, err := http.NewRequest("PUT", baseURL+"/api/trips/"+tripKey+"/update", payload)
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
	getResp, err := http.Get(baseURL + "/api/trips/" + tripKey)
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

// testGetTripProperties retrieves properties for a specific trip
func testGetTripProperties(t *testing.T, baseURL string, tripKey string) {
	url := fmt.Sprintf("%s/api/trips/%s/properties", baseURL, tripKey)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Failed to GET trip properties: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("GET %s returned status %d, want %d. Body: %s", url, resp.StatusCode, http.StatusOK, body)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode trip properties response: %v", err)
	}

	properties, ok := result["properties"].([]interface{})
	if !ok {
		t.Fatal("Expected properties array in response")
	}
	if len(properties) == 0 {
		t.Error("Expected at least one property in trip properties list")
	}

	// Verify properties have expected fields and that at least one is active
	// (trip was created with Hiking and Camping)
	var activeCount int
	for _, p := range properties {
		prop, ok := p.(map[string]interface{})
		if !ok {
			t.Error("Expected property to be an object")
			continue
		}
		if _, ok := prop["name"]; !ok {
			t.Error("Expected property to have a name field")
		}
		if _, ok := prop["active"]; !ok {
			t.Error("Expected property to have an active field")
		}
		if active, _ := prop["active"].(bool); active {
			activeCount++
		}
	}
	if activeCount == 0 {
		t.Error("Expected at least one active property (trip was created with Hiking and Camping)")
	}
}

// testToggleTripProperty toggles a property on the trip and verifies it changed
func testToggleTripProperty(t *testing.T, baseURL string, tripKey string) {
	// Get current properties to find an active one to toggle off
	url := fmt.Sprintf("%s/api/trips/%s/properties", baseURL, tripKey)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Failed to GET trip properties: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode trip properties response: %v", err)
	}

	properties := result["properties"].([]interface{})

	// Find an active property to toggle off
	var targetProp string
	for _, p := range properties {
		prop := p.(map[string]interface{})
		if active, _ := prop["active"].(bool); active {
			targetProp = prop["name"].(string)
			break
		}
	}
	if targetProp == "" {
		t.Skip("No active properties available to toggle")
	}

	// Toggle it off
	toggleURL := fmt.Sprintf("%s/api/trips/%s/properties/%s/toggle", baseURL, tripKey, targetProp)
	toggleResp, err := http.Post(toggleURL, "application/json", nil)
	if err != nil {
		t.Fatalf("Failed to toggle property: %v", err)
	}
	defer toggleResp.Body.Close()

	if toggleResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(toggleResp.Body)
		t.Fatalf("Toggle property returned status %d, want %d. Body: %s", toggleResp.StatusCode, http.StatusOK, body)
	}

	// Verify it is now inactive
	resp2, err := http.Get(url)
	if err != nil {
		t.Fatalf("Failed to GET trip properties after toggle: %v", err)
	}
	defer resp2.Body.Close()

	var result2 map[string]interface{}
	if err := json.NewDecoder(resp2.Body).Decode(&result2); err != nil {
		t.Fatalf("Failed to decode trip properties response after toggle: %v", err)
	}

	properties2 := result2["properties"].([]interface{})
	for _, p := range properties2 {
		prop := p.(map[string]interface{})
		if prop["name"].(string) == targetProp {
			if active, _ := prop["active"].(bool); active {
				t.Errorf("Property %q should be inactive after toggle, but is still active", targetProp)
			}
			return
		}
	}
	t.Errorf("Property %q not found in properties list after toggle", targetProp)
}

// testRenameTrip renames a trip and returns the new trip key
func testRenameTrip(t *testing.T, baseURL string, tripKey string) string {
	newName := "renamed-trip"
	payload := strings.NewReader(fmt.Sprintf(`{"name": %q}`, newName))

	req, err := http.NewRequest("PUT", baseURL+"/api/trips/"+tripKey+"/update", payload)
	if err != nil {
		t.Fatalf("Failed to create rename request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to rename trip: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Rename trip returned status %d, want %d. Body: %s", resp.StatusCode, http.StatusOK, body)
	}

	// Old name should no longer be accessible
	oldResp, err := http.Get(baseURL + "/api/trips/" + tripKey)
	if err != nil {
		t.Fatalf("Failed to GET old trip name: %v", err)
	}
	defer oldResp.Body.Close()
	if oldResp.StatusCode != http.StatusNotFound {
		t.Errorf("Old trip name %q should return 404 after rename, got %d", tripKey, oldResp.StatusCode)
	}

	// New name should be accessible
	newResp, err := http.Get(baseURL + "/api/trips/" + newName)
	if err != nil {
		t.Fatalf("Failed to GET renamed trip: %v", err)
	}
	defer newResp.Body.Close()

	if newResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(newResp.Body)
		t.Fatalf("GET renamed trip returned status %d, want %d. Body: %s", newResp.StatusCode, http.StatusOK, body)
	}

	var trip map[string]interface{}
	if err := json.NewDecoder(newResp.Body).Decode(&trip); err != nil {
		t.Fatalf("Failed to decode renamed trip response: %v", err)
	}

	if name := trip["name"]; name != newName {
		t.Errorf("Expected trip name %q after rename, got %v", newName, name)
	}

	return newName
}
