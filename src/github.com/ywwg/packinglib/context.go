package packinglib

import (
	"fmt"
)

// Property is a value describing a property the context has.  Sort of like boolean flags.
type Property string

// PropertySet is a map holding a list of Properties.  A value of true
// indicates that the Property is allowed.  A value of false indicates
// that the presence of that property is not allowed.  Properties are
// ORed together:  Any allowed Property satisfies the item, but any
// disallowed Property causes the item to reject.
type PropertySet map[Property]bool

var allProperties = PropertySet{
	"Bright":           true,
	"Burn":             true,
	"Camping":          true,
	"Contacts":         true,
	"Dark":             true,
	"Dirt":             true,
	"Fancy":            true,
	"GrumpCamping":     true,
	"HasToiletries":    true,
	"International":    true,
	"Loud":             true,
	"NoCheckedLuggage": true,
	"PaidEvent":        true,
	"Partying":         true,
	"Performing":       true,
	"Sweat":            true,
	"Swimming":         true,
}

func RegisterProperty(prop Property) {
	// Don't worry if the property already exists.
	allProperties[prop] = true
}

func buildPropertySet(allow, disallow []string) PropertySet {
	propSet := make(PropertySet)
	for _, a := range allow {
		if _, ok := allProperties[Property(a)]; !ok {
			panic(fmt.Sprintf("Property not found in allProperties: %s", a))
		}
		propSet[Property(a)] = true
	}
	for _, d := range disallow {
		if _, ok := allProperties[Property(d)]; !ok {
			panic(fmt.Sprintf("Property not found in allProperties: %s", d))
		}
		if _, ok := propSet[Property(d)]; ok {
			panic(fmt.Sprintf("Contradiction: Property already registered in allowed side: %s", d))
		}
		propSet[Property(d)] = false
	}
	return propSet
}

// Context is struct that holds data about the context of the trip
type Context struct {
	// Name of the context ("The Cape", "The Tiny House", "Firefly")
	Name string

	// TemperatureMin is the anticipated minimum temperature.
	TemperatureMin int

	// TemperatureMax is the anticipated maximum temperature.
	TemperatureMax int

	Properties PropertySet
}
