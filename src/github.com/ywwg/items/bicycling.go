package items

import (
	plib "github.com/ywwg/packinglib"
)

var bicycling = []plib.Item{
	plib.NewBasicItem("bike", []string{"Cycling"}, nil),
	plib.NewBasicItem("helmet", []string{"Cycling"}, nil),
	plib.NewBasicItem("bike shoes", []string{"Cycling"}, nil),
	plib.NewBasicItem("bike socks", []string{"Cycling"}, nil),
	plib.NewBasicItem("water bottles", []string{"Cycling"}, nil),
	plib.NewBasicItem("mini bike kit", []string{"Cycling"}, nil),
	plib.NewBasicItem("ankle straps", []string{"Cycling"}, nil),
	plib.NewBasicItem("pump up tires", []string{"Cycling"}, nil),
	plib.NewBasicItem("charge bike lights", []string{"Cycling"}, nil),
	plib.NewBasicItem("GPS map", []string{"LongRide"}, nil),
	plib.NewBasicItem("bike GPS", []string{"LongRide"}, nil),
}

func init() {
	plib.RegisterItems("Bicycling", bicycling)
}
