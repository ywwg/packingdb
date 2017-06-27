package packinglib

import ()

// Property is a value describing a property the context has.  Sort of like boolean flags.
type Property int

const (
	None = iota
	Swimming
	Dirt
	Loud
	Bright
	Dark
	Sweat
	Camping
	Performing
	Burn
	Fancy
	Partying
	NoCheckedLuggage
)

// PropertySet is a map holding a list of Properties.  A value of true
// indicates that the Property is allowed.  A value of false indicates
// that the presence of that property is not allowed.  Properties are
// ORed together:  Any allowed Property satisfies the item, but any
// disallowed Property causes the item to reject.
type PropertySet map[Property]bool

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

var fireflyContext = Context{
	Name:           "Firefly",
	TemperatureMin: 52,
	TemperatureMax: 80,
	Properties: PropertySet{
		Swimming:   true,
		Dirt:       true,
		Loud:       true,
		Bright:     true,
		Sweat:      true,
		Camping:    true,
		Dark:       true,
		Burn:       true,
		Performing: true,
		Partying:   true,
	},
}

func init() {
	RegisterContext("firefly", fireflyContext)
}
