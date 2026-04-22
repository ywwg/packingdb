package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// getTripItems drives GET /api/trips/{name}/items via the real handler stack
// and returns the flattened ItemResponse list for easy lookup by Code.
func getTripItems(t *testing.T, s *Server, tripName string) []ItemResponse {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/trips/%s/items", tripName), nil)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code, "GET items failed: %s", rec.Body.String())

	var resp struct {
		Categories []CategoryResponse `json:"categories"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))

	var out []ItemResponse
	for _, cat := range resp.Categories {
		out = append(out, cat.Items...)
	}
	require.NotEmpty(t, out, "trip %s must have at least one packable item to exercise", tripName)
	return out
}

// TestHTTPEndToEndTripIsolation closes VRFY-02. Drives two trips through the
// real chi router via httptest: creates tripA and tripB, packs an item and
// toggles a property on tripA, then reads tripB back through GET /items and
// GET /properties and asserts zero bleed. Also sanity-checks that tripA's
// own state DID change, so a degenerate "nothing ever mutated" bug cannot pass.
func TestHTTPEndToEndTripIsolation(t *testing.T) {
	s := newTestServer(t)
	createTrip(t, s, "tripA")
	createTrip(t, s, "tripB")

	itemsA := getTripItems(t, s, "tripA")
	itemsB := getTripItems(t, s, "tripB")

	targetCode := itemsA[0].Code

	var bItem *ItemResponse
	for i := range itemsB {
		if itemsB[i].Code == targetCode {
			bItem = &itemsB[i]
			break
		}
	}
	require.NotNil(t, bItem, "tripB must surface the same item code %q as tripA", targetCode)
	require.False(t, itemsA[0].Packed, "precondition: tripA item %q starts unpacked", targetCode)
	require.False(t, bItem.Packed, "precondition: tripB item %q starts unpacked", targetCode)

	{
		req := httptest.NewRequest(http.MethodPost,
			fmt.Sprintf("/api/trips/tripA/items/%s/toggle", targetCode), nil)
		rec := httptest.NewRecorder()
		s.Handler().ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code, "toggle tripA item failed: %s", rec.Body.String())
	}

	{
		req := httptest.NewRequest(http.MethodPost,
			"/api/trips/tripA/properties/outdoors/toggle", nil)
		rec := httptest.NewRecorder()
		s.Handler().ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code, "toggle tripA property failed: %s", rec.Body.String())
	}

	itemsBAfter := getTripItems(t, s, "tripB")
	var bItemAfter *ItemResponse
	for i := range itemsBAfter {
		if itemsBAfter[i].Code == targetCode {
			bItemAfter = &itemsBAfter[i]
			break
		}
	}
	require.NotNil(t, bItemAfter, "tripB must still surface code %q post-mutation", targetCode)
	require.False(t, bItemAfter.Packed, "VRFY-02: tripB item %q must remain unpacked after tripA toggled it", targetCode)

	{
		req := httptest.NewRequest(http.MethodGet, "/api/trips/tripB/properties", nil)
		rec := httptest.NewRecorder()
		s.Handler().ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		var resp struct {
			Properties []PropertyResponse `json:"properties"`
		}
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))

		var found bool
		for _, p := range resp.Properties {
			if p.Name == "outdoors" {
				require.True(t, p.Active,
					"VRFY-02: tripB 'outdoors' must stay Active=true after tripA toggled it off")
				found = true
				break
			}
		}
		require.True(t, found, "tripB properties response must include 'outdoors'")
	}

	itemsAAfter := getTripItems(t, s, "tripA")
	var aItemAfter *ItemResponse
	for i := range itemsAAfter {
		if itemsAAfter[i].Code == targetCode {
			aItemAfter = &itemsAAfter[i]
			break
		}
	}
	require.NotNil(t, aItemAfter, "tripA must still surface code %q post-mutation", targetCode)
	require.True(t, aItemAfter.Packed, "sanity: tripA item %q must be packed after the toggle", targetCode)
}

// TestPackStateSurvivesRestart closes the PROJECT.md Active item:
// "Pack state persists to disk and survives server restart". Flow:
//  1. Build Server s1 on a fresh tempdir and create a trip via the handler.
//  2. Toggle an item packed on s1 via POST /items/{code}/toggle.
//  3. Call s1.persistDirtyTrips() directly (unexported; in-package). No Shutdown
//     needed — StartBackgroundPersist is never called on a test server.
//  4. Discard s1. Construct s2 with NewServer on the SAME dir — scanTrips
//     rebuilds nameToFile from the YAML file on disk.
//  5. GET /api/trips/{name}/items via s2.Handler() and assert the toggled
//     item is still packed=true (proves file-to-memory reload path).
func TestPackStateSurvivesRestart(t *testing.T) {
	dir := t.TempDir()
	r := buildMinimalRegistry(t)

	s1, err := NewServer(r, dir)
	require.NoError(t, err)

	const tripName = "persist-trip"

	{
		body := map[string]interface{}{
			"name":           tripName,
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
		s1.Handler().ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code, "create trip failed: %s", rec.Body.String())
	}

	itemsBefore := getTripItems(t, s1, tripName)
	targetCode := itemsBefore[0].Code

	{
		req := httptest.NewRequest(http.MethodPost,
			fmt.Sprintf("/api/trips/%s/items/%s/toggle", tripName, targetCode), nil)
		rec := httptest.NewRecorder()
		s1.Handler().ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code, "toggle failed: %s", rec.Body.String())
	}

	s1.persistDirtyTrips()

	s2, err := NewServer(r, dir)
	require.NoError(t, err)

	itemsAfter := getTripItems(t, s2, tripName)
	var itemAfter *ItemResponse
	for i := range itemsAfter {
		if itemsAfter[i].Code == targetCode {
			itemAfter = &itemsAfter[i]
			break
		}
	}
	require.NotNil(t, itemAfter, "reloaded trip must surface code %q", targetCode)
	require.True(t, itemAfter.Packed,
		"pack state must survive restart: item %q should be packed after reload", targetCode)
}
