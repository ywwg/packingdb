package items

import (
	plib "github.com/ywwg/packinglib"
)

var bicycling = []plib.Item{
	plib.NewBasicItem("bike", []string{"Cycling"}, nil),
	plib.NewBasicItem("helmet", []string{"Cycling"}, nil),
	plib.NewBasicItem("bike jersey", []string{"Cycling"}, nil),
	plib.NewBasicItem("bike shorts", []string{"Cycling"}, nil),
	plib.NewBasicItem("bike shoes", []string{"Cycling"}, nil),
	plib.NewBasicItem("bike socks", []string{"Cycling"}, nil),
	plib.NewBasicItem("bike gloves", []string{"Cycling"}, nil),
	plib.NewBasicItem("water bottles", []string{"Cycling"}, nil),
	plib.NewBasicItem("mini bike kit", []string{"Cycling"}, nil),
	plib.NewBasicItem("ankle straps", []string{"Cycling"}, nil),
	plib.NewBasicItem("pump up tires", []string{"Cycling"}, nil),
	plib.NewBasicItem("charge bike lights", []string{"Cycling"}, nil),
	plib.NewBasicItem("GPS map", []string{"CyclingLongRide"}, nil),
	plib.NewBasicItem("bike GPS", []string{"CyclingLongRide"}, nil),
}

var climbing = []plib.Item{
	plib.NewBasicItem("climbing shoes", []string{"Climbing"}, nil),
	plib.NewBasicItem("harness", []string{"Climbing"}, nil),
	plib.NewBasicItem("chalk bag", []string{"Climbing"}, nil),
	plib.NewBasicItem("climbing shirts", []string{"Climbing"}, nil),
	plib.NewBasicItem("climbing pants", []string{"Climbing"}, nil),
}

func init() {
	plib.RegisterItems("Bicycling", bicycling)
	plib.RegisterItems("Climbing", climbing)
}
