package contexts

import (
	plib "github.com/ywwg/packinglib"
)

var fireflyContext = plib.Context{
	Name:           "Firefly",
	TemperatureMin: 52,
	TemperatureMax: 80,
	Properties: plib.PropertySet{
		"Swimming":    true,
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

func init() {
	plib.RegisterProperty("Firefly2017")
	plib.RegisterContext("firefly", fireflyContext)
}
