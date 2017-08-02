package contexts

import (
	plib "github.com/ywwg/packinglib"
)

var fireflyContext = plib.Context{
	Name:           "Firefly",
	TemperatureMin: 52,
	TemperatureMax: 80,
	Properties: plib.PropertySet{
		"Dirt":        true,
		"Loud":        true,
		"Bright":      true,
		"Sweat":       true,
		"Camping":     true,
		"Dark":        true,
		"Burn":        true,
		"Performing":  true,
		"Partying":    true,
		"Firefly2017": true,
	},
}

var capeContext = plib.Context{
	Name:           "Cape",
	TemperatureMin: 69,
	TemperatureMax: 120,
	Properties: plib.PropertySet{
		"Bright":        true,
		"Sweat":         true,
		"Swimming":      true,
		"HasToiletries": true,
	},
}

var berlin = plib.Context{
	Name:           "Berlin",
	TemperatureMin: 60,
	TemperatureMax: 120,
	Properties: plib.PropertySet{
		"Bright":        true,
		"Sweat":         true,
		"Swimming":      true,
		"HasToiletries": true,
		"International": true,
		"Partying":      true,
		"Berlin2017":    true,
	},
}

var tinyhouseSummer = plib.Context{
	Name:           "Tiny House Summer",
	TemperatureMin: 60,
	TemperatureMax: 120,
	Properties: plib.PropertySet{
		"Bright":        true,
		"Sweat":         true,
		"HasToiletries": true,
		"Handy":         true,
		"Tiny House":    true,
		"Performing":    true,
	},
}

func init() {
	plib.RegisterProperty("Firefly2017")
	plib.RegisterProperty("Berlin2017")
	plib.RegisterContext(fireflyContext)
	plib.RegisterContext(capeContext)
	plib.RegisterContext(berlin)
	plib.RegisterContext(tinyhouseSummer)
}
