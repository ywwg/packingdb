package contexts

import (
	plib "github.com/ywwg/packingdb/pkg/packinglib"
)

var fireflyContext = plib.Context{
	Name:           "Firefly",
	TemperatureMin: 55,
	TemperatureMax: 95,
	Properties: plib.PropertySet{
		"BYOB":         true,
		"CarRide":      true,
		"Dirt":         true,
		"Loud":         true,
		"Bright":       true,
		"Sweat":        true,
		"Camping":      true,
		"GrumpCamping": true,
		"Dark":         true,
		"Nevermore":    true,
		"Burn":         true,
		"Performing":   true,
		"DJing":        true,
		"Partying":     true,
		"Suiting":      true,
		"Swimming":     true,
		"Speaker":      true,
		"Tarping":      true,
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

var tinyhouseSpring = plib.Context{
	Name:           "Tiny House Spring",
	TemperatureMin: 40,
	TemperatureMax: 120,
	Properties: plib.PropertySet{
		"Art":           true,
		"Bright":        true,
		"CarRide":       true,
		"Dark":          true,
		"Dog":           true,
		"Handy":         true,
		"HasToiletries": true,
		"Hiking":        true,
		"Speaker":       true,
		"Sweat":         true,
		"Tiny House":    true,
	},
}

var tinyhouseSummer = plib.Context{
	Name:           "Tiny House Summer",
	TemperatureMin: 60,
	TemperatureMax: 120,
	Properties: plib.PropertySet{
		"Art":           true,
		"Bright":        true,
		"CarRide":       true,
		"Dark":          true,
		"Dog":           true,
		"Handy":         true,
		"HasToiletries": true,
		"Hiking":        true,
		"Speaker":       true,
		"Sweat":         true,
		"Swimming":      true,
		"Tiny House":    true,
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
		"Hiking":        true,
		"Speaker":       true,
		"Dog":           true,
		"CarRide":       true,
	},
}

var tinyhouseWinter = plib.Context{
	Name:           "Tiny House Winter",
	TemperatureMin: 20,
	TemperatureMax: 50,
	Properties: plib.PropertySet{
		"Bright":        true,
		"CarRide":       true,
		"Dog":           true,
		"Handy":         true,
		"HasToiletries": true,
		"Hiking":        true,
		"Speaker":       true,
		"Tiny House":    true,
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

var tubing = plib.Context{
	Name:           "Tubing",
	TemperatureMin: 60,
	TemperatureMax: 120,
	Properties: plib.PropertySet{
		"Boat":          true,
		"Bright":        true,
		"Camping":       true,
		"CarRide":       true,
		"Dark":          true,
		"Handy":         true,
		"HasToiletries": true,
		"Speaker":       true,
		"Sweat":         true,
		"Swimming":      true,
		"PaidEvent":     true,
	},
}

func Register(r plib.Registry) {
	fancon := fanconContext

	r.RegisterContext(fireflyContext)
	r.RegisterContext(capeContext)
	r.RegisterContext(tinyhouseSpring)
	r.RegisterContext(tinyhouseSummer)
	r.RegisterContext(tinyhouseFall)
	r.RegisterContext(tinyhouseWinter)
	r.RegisterContext(daytrip)
	r.RegisterContext(blanktrip)
	r.RegisterContext(fancon)
	r.RegisterContext(japan)
	r.RegisterContext(quillHill)
	r.RegisterContext(tubing)
}
