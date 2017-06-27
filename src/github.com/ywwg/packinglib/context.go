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
	GrumpCamping
	Performing
	Burn
	Fancy
)

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

var FireflyContext = Context{
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
	},
}
