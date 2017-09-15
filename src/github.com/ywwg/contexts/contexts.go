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
	TemperatureMin: 53,
	TemperatureMax: 80,
	Properties: plib.PropertySet{
		"Bright":        true,
		"Sweat":         true,
		"HasToiletries": true,
		"International": true,
		"Partying":      true,
		"Lodging":       true,
		"Berlin2017":    true,
		"Flight":        true,
		"Fancy":         true,
		"Big Trip":      true,
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
		//"Performing":    true,
	},
}

var sustainRelease = plib.Context{
	Name:           "Sustain-Release",
	TemperatureMin: 56,
	TemperatureMax: 76,
	Properties: plib.PropertySet{
		"Bright":    true,
		"Sweat":     true,
		"Partying":  true,
		"Camping":   true,
		"Dark":      true,
		"Loud":      true,
		"NoFire":    true,
		"Insecure":  true,
		"Swimming":  true,
		"PaidEvent": true,
	},
}

func init() {
	firefly2017 := fireflyContext
	firefly2017.Name = "Firefly2017"

	retreat2017 := capeContext
	retreat2017.Name = "Retreat2017"
	retreat2017.Properties["Retreat2017"] = true

	berlin2017 := berlin
	berlin2017.Name = "Berlin2017"

	plib.RegisterContext(fireflyContext)
	plib.RegisterContext(firefly2017)
	plib.RegisterContext(capeContext)
	plib.RegisterContext(retreat2017)
	plib.RegisterContext(berlin)
	plib.RegisterContext(berlin2017)
	plib.RegisterContext(tinyhouseSummer)
	plib.RegisterContext(sustainRelease)
}
