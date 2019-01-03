package packinglib

import (
	"fmt"
	"sort"
	"strings"
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
	"BigTrip":          true, // this is a big deal
	"Bright":           true, // sun, and therefore sunscreen.
	"Burn":             true,
	"Business":         true, // some work happening.
	"BYOB":             true, // need to pack some booze
	"Camping":          true,
	"Con":              true, // convention of any sort
	"Climbing":         true,
	"Cycling":          true, // moderate amount of biking
	"CyclingLongRide":  true, // biking a long distance
	"Dark":             true, // mostly for camping, but anytime you'll be wandering in the dark
	"Dirt":             true, // are you going to get dirty?
	"DiningOut":        true,
	"Flight":           true,
	"GrumpCamping":     true, // A special type of burn
	"Handy":            true, // need tools?
	"HasToiletries":    true,
	"Insecure":         true, // Don't bring valuables
	"International":    true,
	"Lodging":          true, // Paid lodging like a hotel or airbnb
	"Loud":             true,
	"Nevermore":        true,
	"NoCheckedLuggage": true,
	"NoFire":           true, // Used when there's camping, but no fire allowed at all
	"PaidEvent":        true, // Concerts, shows, cons, etc
	"PaidTravel":       true, // Paying to travel (Flight, rail, bus, etc)
	"Speaker":          true, // Need to play some music
	"Suiting":          true,
	"Partying":         true,
	"Snow":             true,
	"Sweat":            true, // are you gonna sweat up the car?
	"Swimming":         true,
	"Tiny House":       true,
	"UltraFormal":      true, // do you need to *really* dress up?
	// "Performing" for all music, and then add the specific ones
	"Performing": true,
	"Modular":    true,
	"DJing":      true,
}

// RegisterProperty adds a new Property to the database so it can be used.
// Does not verify that all of the properties are in the allProperties map.
func RegisterProperty(prop Property) {
	// Don't worry if the property already exists.
	allProperties[prop] = true
}

func buildPropertySet(allow, disallow []string) PropertySet {
	propSet := make(PropertySet)
	for _, a := range allow {
		propSet[Property(a)] = true
	}
	for _, d := range disallow {
		if _, ok := propSet[Property(d)]; ok {
			panic(fmt.Sprintf("Contradiction: Property already registered in allowed side: %s", d))
		}
		propSet[Property(d)] = false
	}
	return propSet
}

// ListProperties returns all of the registered properties as a slice of strings.
func ListProperties() []string {
	var l []string
	for k := range allProperties {
		l = append(l, string(k))
	}
	less := func(i, j int) bool {
		return strings.ToLower(l[i]) < strings.ToLower(l[j])
	}
	sort.Slice(l, less)
	return l
}
