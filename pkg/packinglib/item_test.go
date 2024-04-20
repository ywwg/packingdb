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

// func TestItemString(t *testing.T) {
// 	// Create a test item
// 	item := &Item{
// 		name:          "Test Item",
// 		count:         2.5,
// 		packed:        true,
// 		prerequisites: make(PropertySet),
// 		mutators:      []PackMutator{},
// 	}

// 	// Call the String method
// 	result := item.String()

// 	// Check if the string representation is correct
// 	expectedString := "● 2.5 Test Item"
// 	if result != expectedString {
// 		t.Errorf("String failed, expected: %s, got: %s", expectedString, result)
// 	}
// }

// // Add more tests for other methods as needed
