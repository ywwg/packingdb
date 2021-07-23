package packinglib

import (
	"fmt"
	"sort"
	"strings"
)

// Property is used to describe an attribute of a context that should be
// user understandable.  Users will choose which properties apply to a given
// context.
type Property string

// PropertySet is a map holding a list of Properties.  A value of true
// indicates that the Property is allowed.  A value of false indicates
// that the presence of that property is not allowed.  Properties are
// ORed together:  Any allowed Property satisfies the item, but any
// disallowed Property causes the item to reject.
type PropertySet map[Property]bool

var allProperties = map[Property]string{
	"Art":              "gonna do some drawing",
	"BigTrip":          "this is a big deal",
	"Boat":             "on a boat! (or floating in a tube)",
	"Bright":           "sun, and therefore sunscreen.",
	"Burn":             "",
	"Business":         "some work happening.",
	"BYOB":             "need to pack some booze",
	"Camping":          "",
	"CarRide":          "includes a long drive",
	"Con":              "convention of any sort",
	"Climbing":         "",
	"Dog":              "Bringing the doggo",
	"LeadClimbing":     "Extra gear for lead",
	"Cycling":          "moderate amount of biking",
	"CyclingLongRide":  "biking a long distance",
	"Dark":             "mostly for camping, but anytime you'll be wandering in the dark",
	"Dirt":             "are you going to get dirty?",
	"DiningOut":        "",
	"Flight":           "",
	"GrumpCamping":     "A special type of burn",
	"Handy":            "need tools?",
	"HasToiletries":    "",
	"Hiking":           "tromping through the woods",
	"Insecure":         "Don't bring valuables",
	"International":    "",
	"Karaoke":          "Belting out some classics",
	"Lodging":          "Paid lodging like a hotel or airbnb",
	"Loud":             "",
	"Nevermore":        "",
	"NoCheckedLuggage": "",
	"NoFire":           "Used when there's camping, but no fire allowed at all",
	"NYC":              "New York City???",
	"PaidEvent":        "Concerts, shows, cons, etc",
	"PaidTravel":       "Paying to travel (Flight, rail, bus, etc)",
	"Projector":        "Movies on the big screen",
	"Speaker":          "Need to play some music",
	"Suiting":          "",
	"Partying":         "",
	"Playtime":         "",
	"Skiing":           "",
	"Snow":             "",
	"Sweat":            "are you gonna sweat up the car?",
	"Swimming":         "",
	"Tarping":          "Putting up tarps",
	"Tiny House":       "",
	"UltraFormal":      "do you need to *really* dress up?",
	"PA System":        "Need the PA",
	"Performing":       "Use this for all music, and then add the specific ones",
	"Modular":          "",
	"DJing":            "",
}

// RegisterProperty adds a new Property to the database so it can be used.
// desc should be a user-visible description of the property.
// Does not verify that all of the properties are in the allProperties map.
func RegisterProperty(prop Property, desc string) {
	// Don't worry if the property already exists.
	allProperties[prop] = desc
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
func ListProperties() []Property {
	var l []Property
	for k := range allProperties {
		l = append(l, k)
	}
	less := func(i, j int) bool {
		return strings.ToLower(string(l[i])) < strings.ToLower(string(l[j]))
	}
	sort.Slice(l, less)
	return l
}
