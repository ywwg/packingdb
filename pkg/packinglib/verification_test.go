package packinglib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestTripIsolationPackAndProperty closes VRFY-01. It creates two trips from
// the same master registry and asserts that BOTH a pack mutation AND a property
// mutation on Trip A leave Trip B and the master registry untouched — the
// single-clone invariant established in Phase 2 must hold across both axes.
//
// Existing TestTripsHaveIsolatedPackState (isolation_test.go) covers only the
// pack axis. This test broadens the guard to the property axis and colocates
// both assertions so a single regression run surfaces isolation breakage on
// either path.
func TestTripIsolationPackAndProperty(t *testing.T) {
	master := NewTestRegistry()

	ctxA := newBeachTripContext(t, master, 3, 70, 90)
	tripA, err := NewTripFromCustomContext(master, ctxA)
	require.NoError(t, err)

	ctxB := newBeachTripContext(t, master, 3, 70, 90)
	tripB, err := NewTripFromCustomContext(master, ctxB)
	require.NoError(t, err)

	tripA.Pack("t-shirt", true)

	tshirtA := findItem(tripA, "t-shirt")
	require.NotNil(t, tshirtA, "tripA: t-shirt must be resolvable via findItem")
	require.True(t, tshirtA.Packed(), "tripA: t-shirt must be packed after Pack(true)")

	tshirtB := findItem(tripB, "t-shirt")
	require.NotNil(t, tshirtB, "tripB: t-shirt must be resolvable via findItem")
	require.False(t, tshirtB.Packed(), "tripB: t-shirt must remain unpacked after tripA.Pack")

	tshirtMaster := findMasterItem(master, "t-shirt")
	require.NotNil(t, tshirtMaster, "master: t-shirt must be resolvable via findMasterItem")
	require.False(t, tshirtMaster.Packed(), "master: t-shirt must remain unpacked after tripA.Pack")

	require.True(t, tripA.HasProperty(Property("outdoors")), "precondition: tripA must start with outdoors=true")
	require.True(t, tripB.HasProperty(Property("outdoors")), "precondition: tripB must start with outdoors=true")

	require.NoError(t, tripA.RemoveProperty("outdoors"))

	require.False(t, tripA.HasProperty(Property("outdoors")),
		"tripA: outdoors must be removed after RemoveProperty")
	require.True(t, tripB.HasProperty(Property("outdoors")),
		"tripB: outdoors must still be active — property mutation on tripA must not bleed")
}
