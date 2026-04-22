package packinglib

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// TestRemovePropertyDeletesKey: HARD-03 Edit 1 regression.
// removeProperty must delete the key, not set it to false (Pitfall 9).
func TestRemovePropertyDeletesKey(t *testing.T) {
	r := NewTestRegistry()
	c, err := NewContext(r, "test-remove-ctx", 3, 50, 80, []string{"outdoors"})
	require.NoError(t, err)

	_, ok := c.Properties[Property("outdoors")]
	require.True(t, ok, "precondition: outdoors should be in the map after NewContext")
	require.True(t, c.hasProperty(Property("outdoors")))

	require.NoError(t, c.removeProperty("outdoors"))

	_, ok = c.Properties[Property("outdoors")]
	require.False(t, ok, "outdoors must be DELETED from Properties map, not set to false")
	require.False(t, c.hasProperty(Property("outdoors")))
}

// TestRemovePropertyEmptyStringNoOp: empty string input is a no-op (preserved).
func TestRemovePropertyEmptyStringNoOp(t *testing.T) {
	r := NewTestRegistry()
	c, err := NewContext(r, "test-empty-ctx", 3, 50, 80, []string{"outdoors"})
	require.NoError(t, err)

	sizeBefore := len(c.Properties)
	require.NoError(t, c.removeProperty(""))
	require.Equal(t, sizeBefore, len(c.Properties))
}

// TestRemovePropertyUnknownReturnsError: removing a never-registered property
// returns an error (behavior preserved from pre-fix code).
func TestRemovePropertyUnknownReturnsError(t *testing.T) {
	r := NewTestRegistry()
	c, err := NewContext(r, "test-unknown-ctx", 3, 50, 80, []string{"outdoors"})
	require.NoError(t, err)

	err = c.removeProperty("never-registered-property")
	require.Error(t, err)
}

// TestFromTripOmitsRemovedProperty: HARD-03 end-to-end.
// Add then remove a property via Trip.RemoveProperty, marshal FromTrip to
// YAML bytes, assert the removed property does not appear in output.
func TestFromTripOmitsRemovedProperty(t *testing.T) {
	r := NewTestRegistry()
	ctx := newBeachTripContext(t, r, 3, 70, 90)
	trip, err := NewTripFromCustomContext(r, ctx)
	require.NoError(t, err)

	require.NoError(t, trip.RemoveProperty("outdoors"))

	yt := FromTrip(trip)
	out, err := yaml.Marshal(yt)
	require.NoError(t, err)

	s := string(out)
	require.NotContains(t, s, "outdoors", "removed property must not appear in YAML output")
	require.Contains(t, s, "hot-weather", "retained property must appear in YAML output")
}

// TestFromTripFiltersStaleFalseEntry: HARD-03 Edit 2 regression.
// Even if a Trip's Context.Properties somehow contains a key with val==false
// (e.g., an in-memory trip loaded before the delete fix migrated its state),
// FromTrip must filter it out — matching the SaveToCSV behavior at trip.go:468.
func TestFromTripFiltersStaleFalseEntry(t *testing.T) {
	r := NewTestRegistry()
	ctx := newBeachTripContext(t, r, 3, 70, 90)
	trip, err := NewTripFromCustomContext(r, ctx)
	require.NoError(t, err)

	trip.C.Properties[Property("stale-flag")] = false
	trip.C.Properties[Property("active-flag")] = true

	yt := FromTrip(trip)
	out, err := yaml.Marshal(yt)
	require.NoError(t, err)

	s := string(out)
	require.False(t, strings.Contains(s, "stale-flag"), "stale false entry must be filtered out of YAML output")
	require.True(t, strings.Contains(s, "active-flag"), "active true entry must appear in YAML output")
}
