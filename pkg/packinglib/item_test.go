package packinglib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var basicContext = Context{
	Name:           "whatever",
	Nights:         3,
	TemperatureMin: 50,
	TemperatureMax: 80,
	Properties:     PropertySet{"prop1": true, "prop2": true},
}

func TestItemAdjustCount(t *testing.T) {
	tests := []struct {
		Name          string
		Prerequisites PropertySet
		Mutators      []packMutator
		Context       Context
		Expected      float64
	}{
		{
			Name:          "No adjustments, comes out as 1",
			Prerequisites: PropertySet{},
			Mutators:      []packMutator{},
			Context:       basicContext,
			Expected:      1.0,
		},
		{
			Name:          "Unfilled requirement, 0",
			Prerequisites: PropertySet{"prop3": true},
			Mutators:      []packMutator{},
			Context:       basicContext,
			Expected:      0.0,
		},
		{
			Name:          "Denial property == 0",
			Prerequisites: PropertySet{"prop1": false},
			Mutators:      []packMutator{},
			Context:       basicContext,
			Expected:      0.0,
		},
		{
			Name:          "Don't have denial, ok",
			Prerequisites: PropertySet{"prop3": false},
			Mutators:      []packMutator{},
			Context:       basicContext,
			Expected:      1.0,
		},
		{
			Name:          "Allowed prop, ok",
			Prerequisites: PropertySet{"prop1": true},
			Mutators:      []packMutator{},
			Context:       basicContext,
			Expected:      1.0,
		},
		{
			Name:          "Temperator Mutator min denies",
			Prerequisites: PropertySet{},
			Mutators:      []packMutator{&temperatureMutator{TemperatureMin: 20, TemperatureMax: 30}},
			Context:       basicContext,
			Expected:      0.0,
		},
		{
			Name:          "Temperator Mutator max denies",
			Prerequisites: PropertySet{},
			Mutators:      []packMutator{&temperatureMutator{TemperatureMin: 90, TemperatureMax: 100}},
			Context:       basicContext,
			Expected:      0.0,
		},
		{
			Name:          "Temperator Mutator allows",
			Prerequisites: PropertySet{},
			Mutators:      []packMutator{&temperatureMutator{TemperatureMin: 30, TemperatureMax: 90}},
			Context:       basicContext,
			Expected:      1.0,
		},
		{
			Name:          "Temperator Mutator allows",
			Prerequisites: PropertySet{},
			Mutators:      []packMutator{&temperatureMutator{TemperatureMin: 30, TemperatureMax: 90}},
			Context:       basicContext,
			Expected:      1.0,
		},
		{
			Name:          "Consumable 3x",
			Prerequisites: PropertySet{},
			Mutators:      []packMutator{&consumableMutator{DailyRate: 3}},
			Context:       basicContext,
			Expected:      9.0,
		},
		{
			Name:          "Max 2",
			Prerequisites: PropertySet{},
			Mutators: []packMutator{
				&consumableMutator{DailyRate: 3},
				&maxCountMutator{Max: 2},
			},
			Context:  basicContext,
			Expected: 2.0,
		},
		{
			Name:          "Custom Func",
			Prerequisites: PropertySet{},
			Mutators: []packMutator{
				&customConsumableMutator{
					RateFunc: func(count float64, nights int, props PropertySet) float64 {
						return count + float64(nights)
					},
				},
			},
			Context:  basicContext,
			Expected: 4.0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			item := Item{
				prerequisites: tc.Prerequisites,
				mutators:      tc.Mutators,
			}

			item.AdjustCount(&tc.Context)
			require.Equal(t, tc.Expected, item.count)
		})
	}
}

func TestItemString(t *testing.T) {
	// Create a test item
	tests := []struct {
		Name     string
		Item     *Item
		Expected string
	}{
		{
			Name:     "empty does not really make sense",
			Item:     &Item{},
			Expected: "",
		},
		{
			Name: "minimal",
			Item: &Item{
				name: "Test Item",
			},
			Expected: "Test Item",
		},
		{
			Name: "simple",
			Item: &Item{
				name:          "Test Item",
				packed:        true,
				prerequisites: make(PropertySet),
				mutators:      []packMutator{},
			},
			Expected: "Test Item",
		},
		{
			Name: "consumable fraction rounds up",
			Item: &Item{
				name:          "consumable",
				units:         "instances",
				packed:        true,
				prerequisites: make(PropertySet),
				mutators: []packMutator{
					&consumableMutator{1.5},
				},
			},
			Expected: "5 instances of consumable",
		},
		{
			Name: "maxconsumable fraction rounds up",
			Item: &Item{
				name:          "consumable",
				units:         "instances",
				packed:        true,
				prerequisites: make(PropertySet),
				mutators: []packMutator{
					&consumableMutator{1},
					&maxCountMutator{2.5},
				},
			},
			Expected: "3 instances of consumable",
		},
		{
			Name: "temperature (noop)",
			Item: &Item{
				name:          "temperature item",
				units:         "instances",
				packed:        true,
				prerequisites: make(PropertySet),
				mutators: []packMutator{
					&consumableMutator{1},
					&temperatureMutator{0, 100},
				},
			},
			Expected: "3 instances of temperature item",
		},
		{
			Name: "consumable nounits",
			Item: &Item{
				name:          "nounits item",
				units:         "",
				packed:        true,
				prerequisites: make(PropertySet),
				mutators: []packMutator{
					&temperatureMutator{0, 100},
					&consumableMutator{2},
				},
			},
			Expected: "6 nounits item",
		},
		{
			Name: "custom consumable",
			Item: &Item{
				name:          "doohickey",
				packed:        true,
				prerequisites: make(PropertySet),
				mutators: []packMutator{
					&customConsumableMutator{func(count float64, nights int, props PropertySet) float64 {
						return 10.0
					}},
				},
			},
			Expected: "10 doohickey",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			tc.Item.AdjustCount(&basicContext)
			got := tc.Item.String()
			require.Equal(t, tc.Expected, got)
		})
	}
}

func TestItemConstruction(t *testing.T) {
	t.Run("basic item", func(t *testing.T) {
		i := NewItem("myitem", []string{"prop1"}, []string{"prop3"})
		i.Units("pounds")
		i.Pack(true)
		i.AdjustCount(&basicContext)
		require.Equal(t, "myitem", i.Name())
		require.Equal(t, "1 pounds of myitem", i.String())
		require.Equal(t, 1.0, i.Count())
		require.Equal(t, true, i.Packed())
		require.Equal(t, PropertySet{"prop1": true, "prop3": false}, i.Prerequisites())
	})

	t.Run("mutators item", func(t *testing.T) {
		i := NewItem("mytempitem", []string{"prop1"}, []string{"prop3"})
		i.TemperatureRange(0, 100)
		i.Consumable(2)
		i.Max(5)
		i.Pack(false)
		i.AdjustCount(&basicContext)
		require.Equal(t, "mytempitem", i.Name())
		require.Equal(t, "5 mytempitem", i.String())
		require.Equal(t, 5.0, i.Count())
		require.Equal(t, false, i.Packed())
		require.Equal(t, PropertySet{"prop1": true, "prop3": false}, i.Prerequisites())
	})

	t.Run("order irrelevant", func(t *testing.T) {
		i := NewItem("mytempitem", []string{"prop1"}, []string{"prop3"})
		i.Max(5)
		i.Pack(false)
		i.Consumable(2)
		i.TemperatureRange(0, 100)
		i.AdjustCount(&basicContext)
		require.Equal(t, "mytempitem", i.Name())
		require.Equal(t, "5 mytempitem", i.String())
		require.Equal(t, 5.0, i.Count())
		require.Equal(t, false, i.Packed())
		require.Equal(t, PropertySet{"prop1": true, "prop3": false}, i.Prerequisites())
	})

	t.Run("custom func", func(t *testing.T) {
		i := NewItem("mycustom", []string{"prop1"}, []string{"prop3"})
		i.Pack(true)
		i.TemperatureRange(0, 100)
		i.Custom(func(count float64, nights int, props PropertySet) float64 {
			return 12.0
		})
		i.AdjustCount(&basicContext)
		require.Equal(t, "mycustom", i.Name())
		require.Equal(t, "12 mycustom", i.String())
		require.Equal(t, 12.0, i.Count())
		require.Equal(t, true, i.Packed())
		require.Equal(t, PropertySet{"prop1": true, "prop3": false}, i.Prerequisites())
	})

}
