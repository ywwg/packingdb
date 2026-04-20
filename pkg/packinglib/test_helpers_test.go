package packinglib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// NewTestRegistry creates a minimal populated registry for use in tests.
// Contains 2 categories with 4 items, 3 properties, and 1 context.
// Item names are distinct from populateRegistry() in yaml_test.go.
func NewTestRegistry() Registry {
	r := NewStructRegistry()

	r.RegisterProperty("hot-weather", "warm climate trip")
	r.RegisterProperty("outdoors", "outdoor activities")
	r.RegisterProperty("formal", "formal occasion")

	r.RegisterItems("clothing", []*Item{
		NewItem("t-shirt", []string{"hot-weather"}, nil),
		NewItem("rain jacket", []string{"outdoors"}, nil),
	})
	r.RegisterItems("gear", []*Item{
		NewItem("flashlight", []string{"outdoors"}, nil),
		NewItem("notebook", nil, nil),
	})

	r.RegisterContext(Context{
		Name:           "beach-trip",
		Nights:         3,
		TemperatureMin: 70,
		TemperatureMax: 90,
		Properties:     PropertySet{"hot-weather": true, "outdoors": true, "beach-trip": true},
	})

	return r
}

func TestTwoRegistriesCoexist(t *testing.T) {
	r1 := NewTestRegistry()
	r2 := NewTestRegistry()

	require.NotNil(t, r1)
	require.NotNil(t, r2)

	// Both registries have the same item categories, independently registered
	require.Equal(t, len(r1.AllItems()), len(r2.AllItems()))

	// Both registries have the same properties
	require.Equal(t, len(r1.AllProperties()), len(r2.AllProperties()))

	// Verify each registry has the expected context
	c1, err := r1.GetContext("beach-trip")
	require.NoError(t, err)
	require.Equal(t, 3, c1.Nights)

	c2, err := r2.GetContext("beach-trip")
	require.NoError(t, err)
	require.Equal(t, 3, c2.Nights)
}

func TestNewContextIdempotent(t *testing.T) {
	r := NewTestRegistry()

	// Create a new context -- calls RegisterProperty and RegisterContext internally
	c1, err := NewContext(r, "test-ctx", 3, 50, 80, []string{"outdoors"})
	require.NoError(t, err)
	require.NotNil(t, c1)

	// Create the same context again -- must NOT panic
	c2, err := NewContext(r, "test-ctx", 5, 40, 70, []string{"outdoors"})
	require.NoError(t, err)
	require.NotNil(t, c2)

	// First registration wins in the registry
	stored, err := r.GetContext("test-ctx")
	require.NoError(t, err)
	require.Equal(t, 3, stored.Nights, "first registration should be preserved")
}

func TestRegistryCountTwo(t *testing.T) {
	// Simulate what -count=2 does: register the same items on fresh registries
	r1 := NewStructRegistry()
	r1.RegisterProperty("Business", "business cat")
	r1.RegisterProperty("Flight", "wheee")
	r1.RegisterProperty("International", "le whee")
	r1.RegisterProperty("camping", "tent")
	r1.RegisterItems("business", []*Item{
		NewItem("sim-laptop", []string{"Business"}, nil),
		NewItem("sim-socks", []string{"Flight"}, []string{"camping"}),
	})

	// Second registry with the same structure should not panic
	r2 := NewStructRegistry()
	r2.RegisterProperty("Business", "business cat")
	r2.RegisterProperty("Flight", "wheee")
	r2.RegisterProperty("International", "le whee")
	r2.RegisterProperty("camping", "tent")
	r2.RegisterItems("business", []*Item{
		NewItem("sim-laptop", []string{"Business"}, nil),
		NewItem("sim-socks", []string{"Flight"}, []string{"camping"}),
	})

	require.Equal(t, len(r1.AllItems()), len(r2.AllItems()))
}
