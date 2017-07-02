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
	TemperatureMin: 60,
	TemperatureMax: 120,
	Properties: plib.PropertySet{
		"Bright":   true,
		"Sweat":    true,
		"Swimming": true,
	},
}

func init() {
	plib.RegisterProperty("Firefly2017")
	plib.RegisterContext("firefly", fireflyContext)
	plib.RegisterContext("Cape", capeContext)
}
