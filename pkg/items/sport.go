package items

import (
	plib "github.com/ywwg/packingdb/pkg/packinglib"
)

var bicycling = []*plib.Item{
	plib.NewItem("bike", []string{"Cycling"}, nil),
	plib.NewItem("helmet", []string{"Cycling"}, nil),
	plib.NewItem("bike jersey", []string{"Cycling"}, nil),
	plib.NewItem("bike shorts", []string{"Cycling"}, nil),
	plib.NewItem("bike shoes", []string{"Cycling"}, nil),
	plib.NewItem("bike socks", []string{"Cycling"}, nil),
	plib.NewItem("bike gloves", []string{"Cycling"}, nil),
	plib.NewItem("water bottles", []string{"Cycling"}, nil),
	plib.NewItem("mini bike kit", []string{"Cycling"}, nil),
	plib.NewItem("ankle straps", []string{"Cycling"}, nil),
	plib.NewItem("pump up tires", []string{"Cycling"}, nil),
	plib.NewItem("charge bike lights", []string{"Cycling"}, nil),
	plib.NewItem("GPS map", []string{"CyclingLongRide"}, nil),
	plib.NewItem("bike GPS", []string{"CyclingLongRide"}, nil),
	plib.NewItem("GPS Dongle", []string{"CyclingLongRide"}, nil),
}

var climbing = []*plib.Item{
	plib.NewItem("climbing shoes", []string{"Climbing"}, nil),
	plib.NewItem("climbing harness", []string{"Climbing"}, nil),
	plib.NewItem("chalk bag", []string{"Climbing"}, nil),
	plib.NewItem("climbing shirts", []string{"Climbing"}, nil),
	plib.NewItem("climbing pants", []string{"Climbing"}, nil),
	plib.NewItem("climbing helmet", []string{"LeadClimbing"}, nil),
	plib.NewItem("belay device", []string{"LeadClimbing"}, nil),
}

var skiing = []*plib.Item{
	plib.NewItem("goggles", []string{"Skiing"}, nil),
	plib.NewItem("warm balaklava", []string{"Skiing"}, nil),
}

var boating = []*plib.Item{
	plib.NewItem("phone dry bag", []string{"Boat"}, nil),
}
