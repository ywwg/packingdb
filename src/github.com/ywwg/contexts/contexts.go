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
		"Flight":        true,
		"DiningOut":     true,
		"BigTrip":       true,
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
		"Dog":           true,
		"CarRide":       true,
		"Projector":     true,
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
		"Dog":           true,
		"CarRide":       true,
		"Projector":     true,
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
		"Dog":           true,
		"CarRide":       true,
		"Projector":     true,
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

var paxUnpluggedContext = plib.Context{
	Name:           "PAX Unplugged",
	TemperatureMin: 35,
	TemperatureMax: 56,
	Properties: plib.PropertySet{
		"Con":           true,
		"Lodging":       true,
		"PaidEvent":     true,
		"PaidTravel":    true,
		"HasToiletries": true,
	},
}

var daytrip = plib.Context{
	Name:           "daytrip",
	TemperatureMin: 25,
	TemperatureMax: 55,
	Properties: plib.PropertySet{
		"Lodging": true,
		"Loud":    true,
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
		"DiningOut":     true,
	},
}

var fanconContext = plib.Context{
	Name:           "fancon",
	TemperatureMin: 43,
	TemperatureMax: 60,
	Properties: plib.PropertySet{
		"Con":        true,
		"Bright":     true,
		"Partying":   true,
		"Loud":       true,
		"Swimming":   true,
		"PaidEvent":  true,
		"Lodging":    true,
		"Performing": true,
		"DJing":      true,
		"Art":        true,
		"BYOB":       true,
		"Speaker":    true,
		"Suiting":    true,
	},
}

var japan = plib.Context{
	Name:           "Japan",
	TemperatureMin: 63,
	TemperatureMax: 80,
	Properties: plib.PropertySet{
		"BigTrip":       true,
		"Bright":        true,
		"Business":      true,
		"DiningOut":     true,
		"Flight":        true,
		"HasToiletries": true,
		"International": true,
		"Lodging":       true,
		"Loud":          true,
		"PaidEvent":     true,
		"PaidTravel":    true,
		"Swimming":      true,
		"Partying":      true,
		"Sweat":         true,
	},
}

var quillHill = plib.Context{
	Name:           "Quill Hill",
	TemperatureMin: 60,
	TemperatureMax: 120,
	Properties: plib.PropertySet{
		"Tiny House Summer": true,
		"Projector":         true,
		"PA System":         true,
		"Karaoke":           true,
	},
}

func init() {
	nectr2017 := nectrContext
	nectr2017.Name = "Nectr2017"

	fancon := fanconContext

	plib.RegisterContext(fireflyContext)
	plib.RegisterContext(capeContext)
	plib.RegisterContext(tinyhouseSummer)
	plib.RegisterContext(tinyhouseFall)
	plib.RegisterContext(tinyhouseWinter)
	plib.RegisterContext(daytrip)
	plib.RegisterContext(blanktrip)
	plib.RegisterContext(fancon)
	plib.RegisterContext(japan)
	plib.RegisterContext(quillHill)
}
