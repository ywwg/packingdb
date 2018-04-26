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
	"Art":              true, // gonna do some drawing
	"Big Trip":         true, // this is a big deal
	"Bright":           true, // sun!
	"Burn":             true,
	"Camping":          true,
	"Cycling":          true,
	"Dark":             true, // mostly for camping, but anytime you'll be wandering in the dark
	"Dirt":             true, // are you going to get dirty?
	"Fancy":            true,
	"Flight":           true,
	"Formal":           true, // do you need to *really* dress up?
	"GrumpCamping":     true,
	"Handy":            true, // need tools?
	"HasToiletries":    true,
	"Insecure":         true, // Don't bring valuables
	"International":    true,
	"Lodging":          true, // Paid lodging like a hotel or airbnb
	"Loud":             true,
	"NoCheckedLuggage": true,
	"NoFire":           true, // Used when there's camping, but no fire allowed at all
	"PaidEvent":        true,
	"Partying":         true,
	"Sweat":            true, // are you gonna sweat up the car?
	"Swimming":         true,
	"Tiny House":       true,
	// "Performing" for all music, and then add the specific ones
	"Performing": true,
	"Modular":    true,
	"DJing":      true,
}

// RegisterProperty adds a new Property to the database so it can be used.
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
