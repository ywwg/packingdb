package contexts

import (
	plib "github.com/ywwg/packinglib"
)

var fireflyContext = plib.Context{
	Name:           "Firefly",
	TemperatureMin: 55,
	TemperatureMax: 95,
	Properties: plib.PropertySet{
		"Dirt":         true,
		"Loud":         true,
		"Bright":       true,
		"Sweat":        true,
		"Camping":      true,
		"GrumpCamping": true,
		"Dark":         true,
		"Burn":         true,
		"Performing":   true,
		"Modular":      true,
		"DJing":        true,
		"Partying":     true,
		"Firefly2018":  true,
		"Suiting":      true,
		"Swimming":     true,
	},
}

var nectrContext = plib.Context{
	Name:           "Nectr",
	TemperatureMin: 40,
	TemperatureMax: 80,
	Properties: plib.PropertySet{
		"Dirt":         true,
		"Loud":         true,
		"Bright":       true,
		"Sweat":        true,
		"Camping":      true,
		"Dark":         true,
		"Burn":         true,
		"Performing":   true,
		"Partying":     true,
		"Speaker":      true,
		"GrumpCamping": true,
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
		"Dark":          true,
		"Art":           true,
		"Speaker":       true,
	},
}

var tinyhouseFall = plib.Context{
	Name:           "Tiny House Fall",
	TemperatureMin: 50,
	TemperatureMax: 120,
	Properties: plib.PropertySet{
		"Bright":        true,
		"Sweat":         true,
		"HasToiletries": true,
		"Handy":         true,
		"Tiny House":    true,
		"Speaker":       true,
	},
}

var tinyhouseWinter = plib.Context{
	Name:           "Tiny House Winter",
	TemperatureMin: 20,
	TemperatureMax: 50,
	Properties: plib.PropertySet{
		"Bright":        true,
		"HasToiletries": true,
		"Handy":         true,
		"Tiny House":    true,
		"Speaker":       true,
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
		"Speaker":   true,
	},
}

var offsiteContext = plib.Context{
	Name:           "Offsite",
	TemperatureMin: 65,
	TemperatureMax: 120,
	Properties: plib.PropertySet{
		"Bright":        true,
		"Sweat":         true,
		"Swimming":      true,
		"HasToiletries": true,
		"Lodging":       true,
	},
}

var paxUnpluggedContext = plib.Context{
	Name:           "PAX Unplugged",
	TemperatureMin: 35,
	TemperatureMax: 56,
	Properties: plib.PropertySet{
		"Lodging":   true,
		"PaidEvent": true,
	},
}

var daytrip = plib.Context{
	Name:           "daytrip",
	TemperatureMin: 25,
	TemperatureMax: 55,
	Properties: plib.PropertySet{
		"Lodging":   true,
		"PaidEvent": true,
		"Loud":      true,
	},
}

var blanktrip = plib.Context{
	Name:           "blanktrip",
	TemperatureMin: 0,
	TemperatureMax: 120,
	Properties:     plib.PropertySet{},
}

var florida = plib.Context{
	Name:           "Florida",
	TemperatureMin: 64,
	TemperatureMax: 82,
	Properties: plib.PropertySet{
		"Bright":        true,
		"Sweat":         true,
		"Swimming":      true,
		"HasToiletries": true,
		"Lodging":       true,
		"Flight":        true,
		"Fancy":         true,
	},
}

var conContext = plib.Context{
	Name:           "Con",
	TemperatureMin: 43,
	TemperatureMax: 60,
	Properties: plib.PropertySet{
		"Bright":     true,
		"Partying":   true,
		"Loud":       true,
		"Swimming":   true,
		"PaidEvent":  true,
		"Lodging":    true,
		"Performing": true,
		"DJing":      true,
		"Con":        true,
		"Formal":     true,
		"Art":        true,
		"Drinking":   true,
		"Speaker":    true,
	},
}

func init() {
	firefly2017 := fireflyContext
	firefly2017.Name = "Firefly2017"

	firefly2018 := fireflyContext
	firefly2018.Name = "Firefly2018"

	// retreat2017 := capeContext
	// retreat2017.Name = "Retreat2017"
	// retreat2017.Properties["Retreat2017"] = true

	berlin2017 := berlin
	berlin2017.Name = "Berlin2017"

	offsite2017 := offsiteContext

	nectr2017 := nectrContext
	nectr2017.Name = "Nectr2017"

	con := conContext

	plib.RegisterContext(fireflyContext)
	plib.RegisterContext(firefly2017)
	plib.RegisterContext(firefly2018)
	plib.RegisterContext(capeContext)
	// plib.RegisterContext(retreat2017)
	plib.RegisterContext(berlin)
	plib.RegisterContext(berlin2017)
	plib.RegisterContext(tinyhouseSummer)
	plib.RegisterContext(tinyhouseFall)
	plib.RegisterContext(tinyhouseWinter)
	plib.RegisterContext(sustainRelease)
	plib.RegisterContext(offsite2017)
	plib.RegisterContext(nectr2017)
	plib.RegisterContext(paxUnpluggedContext)
	plib.RegisterContext(daytrip)
	plib.RegisterContext(blanktrip)
	plib.RegisterContext(florida)
	plib.RegisterContext(con)
}
