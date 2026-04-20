package packinglib

import (
	"errors"
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

	if err := r.RegisterContext(Context{
		Name:           "beach-trip",
		Nights:         3,
		TemperatureMin: 70,
		TemperatureMax: 90,
		Properties:     PropertySet{"hot-weather": true, "outdoors": true, "beach-trip": true},
	}); err != nil {
		panic(err)
	}

	return r
}

func TestTwoRegistriesCoexist(t *testing.T) {
	r1 := NewTestRegistry()
	r2 := NewTestRegistry()

	require.NotNil(t, r1)
	require.NotNil(t, r2)

	// Both registries have the same set of categories
	cats1 := make([]Category, 0, len(r1.AllItems()))
	for c := range r1.AllItems() {
		cats1 = append(cats1, c)
	}
	cats2 := make([]Category, 0, len(r2.AllItems()))
	for c := range r2.AllItems() {
		cats2 = append(cats2, c)
	}
	require.ElementsMatch(t, cats1, cats2)

	// Both registries have the same set of items (by name) in each category
	for _, cat := range cats1 {
		names1 := make([]string, 0, len(r1.AllItems()[cat]))
		for _, i := range r1.AllItems()[cat] {
			names1 = append(names1, i.Name())
		}
		names2 := make([]string, 0, len(r2.AllItems()[cat]))
		for _, i := range r2.AllItems()[cat] {
			names2 = append(names2, i.Name())
		}
		require.ElementsMatch(t, names1, names2, "items in category %s", cat)
	}

	// Both registries have the same set of properties (by key)
	props1 := make([]Property, 0, len(r1.AllProperties()))
	for p := range r1.AllProperties() {
		props1 = append(props1, p)
	}
	props2 := make([]Property, 0, len(r2.AllProperties()))
	for p := range r2.AllProperties() {
		props2 = append(props2, p)
	}
	require.ElementsMatch(t, props1, props2)

	// Verify each registry has the expected context
	c1, err := r1.GetContext("beach-trip")
	require.NoError(t, err)
	require.Equal(t, 3, c1.Nights)

	c2, err := r2.GetContext("beach-trip")
	require.NoError(t, err)
	require.Equal(t, 3, c2.Nights)
}

func TestNewContextDuplicateReturnsError(t *testing.T) {
	r := NewTestRegistry()

	// Create a new context -- calls RegisterProperty and RegisterContext internally
	c1, err := NewContext(r, "test-ctx", 3, 50, 80, []string{"outdoors"})
	require.NoError(t, err)
	require.NotNil(t, c1)

	// Duplicate registration returns ErrContextExists instead of silently succeeding
	c2, err := NewContext(r, "test-ctx", 5, 40, 70, []string{"outdoors"})
	require.ErrorIs(t, err, ErrContextExists)
	require.Nil(t, c2)

	// Original registration is preserved in the registry
	stored, err := r.GetContext("test-ctx")
	require.NoError(t, err)
	require.Equal(t, 3, stored.Nights, "first registration should be preserved")
}

func TestRegisterContextDuplicateReturnsError(t *testing.T) {
	r := NewStructRegistry()
	c := Context{
		Name:           "beach",
		Nights:         3,
		TemperatureMin: 70,
		TemperatureMax: 90,
		Properties:     PropertySet{},
	}

	require.NoError(t, r.RegisterContext(c))

	err := r.RegisterContext(c)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrContextExists))
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

	names1 := make([]string, 0)
	for _, items := range r1.AllItems() {
		for _, i := range items {
			names1 = append(names1, i.Name())
		}
	}
	names2 := make([]string, 0)
	for _, items := range r2.AllItems() {
		for _, i := range items {
			names2 = append(names2, i.Name())
		}
	}
	require.ElementsMatch(t, names1, names2)
}
