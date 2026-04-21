package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ywwg/packingdb/pkg/packinglib"
)

// buildMinimalRegistry mirrors the shape of packinglib.NewTestRegistry (which
// lives in a _test.go file in the packinglib package and is therefore not
// importable here). Keep this tiny — the race test does not care about the
// item catalog, only that there is one trip with at least one togglable item.
func buildMinimalRegistry(t *testing.T) packinglib.Registry {
	t.Helper()
	r := packinglib.NewStructRegistry()
	r.RegisterProperty(packinglib.Property("hot-weather"), "warm climate trip")
	r.RegisterProperty(packinglib.Property("outdoors"), "outdoor activities")
	r.RegisterItems(packinglib.Category("clothing"), []*packinglib.Item{
		packinglib.NewItem("t-shirt", []string{"hot-weather"}, nil),
		packinglib.NewItem("rain jacket", []string{"outdoors"}, nil),
	})
	r.RegisterItems(packinglib.Category("gear"), []*packinglib.Item{
		packinglib.NewItem("flashlight", []string{"outdoors"}, nil),
		packinglib.NewItem("notebook", nil, nil),
	})
	return r
}

// newTestServer constructs a Server wired to a temp trips dir. t.TempDir
// handles cleanup automatically.
func newTestServer(t *testing.T) *Server {
	t.Helper()
	dir := t.TempDir()

	r := buildMinimalRegistry(t)
	s, err := NewServer(r, dir)
	require.NoError(t, err)
	return s
}

// createTrip drives POST /api/trips through the real handler stack so the
// trip ends up in the same cache the read/write handlers will hit.
func createTrip(t *testing.T, s *Server, name string) {
	t.Helper()
	body := map[string]interface{}{
		"name":           name,
		"nights":         3,
		"temperatureMin": 70,
		"temperatureMax": 90,
		"properties":     []string{"hot-weather", "outdoors"},
	}
	buf, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/trips", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code, "create trip failed: %s", rec.Body.String())
}

// TestReadHandlersRaceSafeAgainstToggle is the HARD-02 regression test.
// It spawns a reader goroutine (hammering the three read handlers) and a
// writer goroutine (hammering the toggle handler) against the same trip
// and relies on `go test -race` to catch any data race. The test must run
// cleanly under `-race` with no race reports and no panics.
//
// Before HARD-02 fix: getTripHandler / getItemsHandler / getTripPropertiesHandler
// read trip.C.Properties and trip.packList without holding s.mu, so the
// race detector flags a conflict with the toggleItemHandler writer.
// After the fix: reads hold s.mu.RLock for their duration and toggle holds
// s.mu.Lock, so the two serialize cleanly.
func TestReadHandlersRaceSafeAgainstToggle(t *testing.T) {
	const iterations = 200

	s := newTestServer(t)
	const tripName = "race-test-trip"
	createTrip(t, s, tripName)

	handler := s.Handler()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		readPaths := []string{
			fmt.Sprintf("/api/trips/%s", tripName),
			fmt.Sprintf("/api/trips/%s/items", tripName),
			fmt.Sprintf("/api/trips/%s/properties", tripName),
		}
		for i := 0; i < iterations; i++ {
			path := readPaths[i%len(readPaths)]
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			if rec.Code >= 500 {
				t.Errorf("reader got 5xx on %s: %d %s", path, rec.Code, rec.Body.String())
				return
			}
		}
	}()

	// Writer: toggles the same item repeatedly. The code "a" may or may not
	// resolve to a real item for this registry; even a 400 still drives the
	// writer handler through its s.mu.Lock path concurrently with the readers,
	// which is all the race detector needs to observe.
	go func() {
		defer wg.Done()
		path := fmt.Sprintf("/api/trips/%s/items/a/toggle", tripName)
		for i := 0; i < iterations; i++ {
			req := httptest.NewRequest(http.MethodPost, path, nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			if rec.Code >= 500 {
				t.Errorf("writer got 5xx on %s: %d %s", path, rec.Code, rec.Body.String())
				return
			}
		}
	}()

	wg.Wait()
}
