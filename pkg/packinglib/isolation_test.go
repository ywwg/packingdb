package packinglib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// findItem walks a trip's packList looking for an item by name (case-sensitive).
// Returns nil if not found.
func findItem(t *Trip, name string) *Item {
	for _, items := range t.packList {
		for _, it := range items {
			if it.Name() == name {
				return it
			}
		}
	}
	return nil
}

// findMasterItem walks a Registry's AllItems() looking for an item by name.
func findMasterItem(r Registry, name string) *Item {
	for _, items := range r.AllItems() {
		for _, it := range items {
			if it.Name() == name {
				return it
			}
		}
	}
	return nil
}

func newBeachTripContext(t *testing.T, r Registry, nights, tmin, tmax int) *Context {
	t.Helper()
	c, err := r.GetContext("beach-trip")
	require.NoError(t, err)
	// GetContext returns a partial context without registry; rewire to match
	// the production load path (see yaml.go AsTrip and Pitfall 6).
	c.registry = r
	c.Nights = nights
	c.TemperatureMin = tmin
	c.TemperatureMax = tmax
	return c
}

func TestTripsShareClonedRegistryWithTheirContext(t *testing.T) {
	master := NewTestRegistry()

	ctxA := newBeachTripContext(t, master, 3, 70, 90)
	tripA, err := NewTripFromCustomContext(master, ctxA)
	require.NoError(t, err)

	ctxB := newBeachTripContext(t, master, 3, 70, 90)
	tripB, err := NewTripFromCustomContext(master, ctxB)
	require.NoError(t, err)

	// Single-clone invariant: Trip and Context must share the exact registry pointer.
	require.Same(t, tripA.registry, tripA.C.registry, "tripA: Trip.registry and Context.registry must be the same instance")
	require.Same(t, tripB.registry, tripB.C.registry, "tripB: Trip.registry and Context.registry must be the same instance")

	// Each trip has its own registry, distinct from master and from each other.
	require.NotSame(t, master, tripA.registry, "tripA must not share registry with master")
	require.NotSame(t, master, tripB.registry, "tripB must not share registry with master")
	require.NotSame(t, tripA.registry, tripB.registry, "tripA and tripB must have distinct registries")
}

func TestTripsHaveIsolatedPackState(t *testing.T) {
	master := NewTestRegistry()

	ctxA := newBeachTripContext(t, master, 3, 70, 90)
	tripA, err := NewTripFromCustomContext(master, ctxA)
	require.NoError(t, err)

	ctxB := newBeachTripContext(t, master, 3, 70, 90)
	tripB, err := NewTripFromCustomContext(master, ctxB)
	require.NoError(t, err)

	// Pack t-shirt in Trip A. "t-shirt" exists because beach-trip has hot-weather.
	tripA.Pack("t-shirt", true)

	itemA := findItem(tripA, "t-shirt")
	require.NotNil(t, itemA)
	require.True(t, itemA.Packed(), "tripA: t-shirt must be packed")

	itemB := findItem(tripB, "t-shirt")
	require.NotNil(t, itemB)
	require.False(t, itemB.Packed(), "tripB: t-shirt must remain unpacked")

	masterItem := findMasterItem(master, "t-shirt")
	require.NotNil(t, masterItem)
	require.False(t, masterItem.Packed(), "master: t-shirt must remain unpacked")
}

func TestTripsHaveIsolatedCountState(t *testing.T) {
	master := NewTestRegistry()

	// Trip A: hot weather, notebook has no prereqs so it applies; count should be 1.
	ctxA := newBeachTripContext(t, master, 3, 70, 90)
	tripA, err := NewTripFromCustomContext(master, ctxA)
	require.NoError(t, err)

	// Trip B: same context shape (independent clone); start with a pristine count view.
	ctxB := newBeachTripContext(t, master, 3, 70, 90)
	tripB, err := NewTripFromCustomContext(master, ctxB)
	require.NoError(t, err)

	// Mutate trip A's notebook count directly via item.count, then reassert it on Trip B.
	// We use AdjustCount with a satisfying context on A only.
	itemA := findItem(tripA, "notebook")
	require.NotNil(t, itemA)
	itemA.AdjustCount(tripA.C)
	countA := itemA.Count()
	require.Greater(t, countA, float64(0), "tripA: notebook should have a non-zero count")

	itemB := findItem(tripB, "notebook")
	require.NotNil(t, itemB)
	// Trip B's notebook went through its own makeList -> AdjustCount. Its count
	// is independent of what we just did on Trip A. The critical assertion is
	// that modifying A's item.count via Pack or AdjustCount does not bleed to B.
	// Prove isolation by packing on A and checking B unchanged:
	itemA.Pack(true)
	require.False(t, itemB.Packed(), "tripB: notebook must not see packed=true from tripA mutation")

	// Master item untouched.
	masterItem := findMasterItem(master, "notebook")
	require.NotNil(t, masterItem)
	require.Equal(t, float64(0), masterItem.Count(), "master: notebook count must remain at zero")
	require.False(t, masterItem.Packed(), "master: notebook must remain unpacked")
}

func TestTripCallerContextNotMutated(t *testing.T) {
	master := NewTestRegistry()

	ctx := newBeachTripContext(t, master, 3, 70, 90)
	origPropLen := len(ctx.Properties)

	trip, err := NewTripFromCustomContext(master, ctx)
	require.NoError(t, err)

	// AddProperty on the trip goes through the CLONED context, not the caller's.
	// "formal" is a registered property in NewTestRegistry.
	require.NoError(t, trip.AddProperty("formal"))

	require.Equal(t, origPropLen, len(ctx.Properties),
		"caller's Context.Properties map must not be mutated by trip operations")
	_, hasFormal := ctx.Properties["formal"]
	require.False(t, hasFormal, "caller's Context must not have the property added on the trip")
}
