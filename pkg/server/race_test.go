package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestMultiTripConcurrentPackRace closes VRFY-03. Four goroutines hammer
// toggle requests against four DIFFERENT trips through the real chi router,
// under `go test -race`. Proves two things at once:
//
//  1. Race-safety: the global sync.RWMutex pairing in toggleItemHandler is
//     safe across concurrent writers on distinct trips.
//  2. Cross-trip isolation under load: trip A's toggle never flips trip B/C/D's
//     packed bit, even with 4 writers racing. Post-race, each trip's target
//     item's packed state matches its own per-goroutine successful-toggle parity.
//
// Relation to Phase 3's TestReadHandlersRaceSafeAgainstToggle:
// Phase 3 proved reader/writer pairing on the SAME trip. VRFY-03 proves
// writer/writer pairing across DIFFERENT trips plus the Phase 2 clone invariant
// holding under concurrent load.
func TestMultiTripConcurrentPackRace(t *testing.T) {
	const (
		numTrips   = 4
		iterations = 200
	)

	s := newTestServer(t)
	handler := s.Handler()

	tripNames := []string{"tripA", "tripB", "tripC", "tripD"}
	require.Len(t, tripNames, numTrips)

	for _, name := range tripNames {
		createTrip(t, s, name)
	}

	targetCodes := make([]string, numTrips)
	for i, name := range tripNames {
		req := httptest.NewRequest(http.MethodGet,
			fmt.Sprintf("/api/trips/%s/items", name), nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code,
			"pre-race GET items for %s failed: %s", name, rec.Body.String())

		var resp struct {
			Categories []CategoryResponse `json:"categories"`
		}
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))

		var code string
	outer:
		for _, cat := range resp.Categories {
			for _, it := range cat.Items {
				code = it.Code
				break outer
			}
		}
		require.NotEmpty(t, code, "trip %s must expose at least one packable item code", name)
		targetCodes[i] = code
	}

	successes := make([]int, numTrips)

	var wg sync.WaitGroup
	wg.Add(numTrips)

	for i := 0; i < numTrips; i++ {
		i := i
		go func() {
			defer wg.Done()
			path := fmt.Sprintf("/api/trips/%s/items/%s/toggle",
				tripNames[i], targetCodes[i])
			for n := 0; n < iterations; n++ {
				req := httptest.NewRequest(http.MethodPost, path, nil)
				rec := httptest.NewRecorder()
				handler.ServeHTTP(rec, req)
				if rec.Code >= 500 {
					t.Errorf("goroutine %d (%s): writer got 5xx: %d %s",
						i, tripNames[i], rec.Code, rec.Body.String())
					return
				}
				if rec.Code == http.StatusOK {
					successes[i]++
				}
			}
		}()
	}

	wg.Wait()

	for i, name := range tripNames {
		req := httptest.NewRequest(http.MethodGet,
			fmt.Sprintf("/api/trips/%s/items", name), nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code,
			"post-race GET items for %s failed: %s", name, rec.Body.String())

		var resp struct {
			Categories []CategoryResponse `json:"categories"`
		}
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))

		var found *ItemResponse
		for ci := range resp.Categories {
			for ii := range resp.Categories[ci].Items {
				if resp.Categories[ci].Items[ii].Code == targetCodes[i] {
					found = &resp.Categories[ci].Items[ii]
					break
				}
			}
			if found != nil {
				break
			}
		}
		require.NotNilf(t, found, "post-race: trip %s must still surface code %q",
			name, targetCodes[i])

		expected := successes[i]%2 == 1
		require.Equalf(t, expected, found.Packed,
			"VRFY-03: trip %s item %q should be packed=%v after %d successful toggles",
			name, targetCodes[i], expected, successes[i])

		require.Greaterf(t, successes[i], 0,
			"goroutine %d produced zero successful toggles; parity check would be trivial",
			i)
	}
}
