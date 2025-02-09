package items

import (
	plib "github.com/ywwg/packinglib"
)

var campStuff = []plib.Item{
	plib.NewBasicItem("windshield sun shade", []string{"Camping"}, nil),
	plib.NewBasicItem("tent", []string{"Camping"}, nil),
	plib.NewBasicItem("sleeping bag", []string{"Camping"}, nil),
	plib.NewBasicItem("camping cot", []string{"Camping"}, nil),
	plib.NewTemperatureItem("wool blanket", 0, 44, []string{"Camping"}, nil),
	plib.NewBasicItem("2 pillow cases", []string{"Camping"}, nil),
	plib.NewBasicItem("tent light", []string{"Camping"}, nil),
	plib.NewBasicItem("headlamp", []string{"Dark", "Hiking"}, nil),
	plib.NewBasicItem("towel for butt driving home", []string{"Burn"}, nil),
	plib.NewBasicItem("tarps", []string{"Burn", "Quill Hill", "Tarping"}, []string{"Firefly"}),
	plib.NewBasicItem("pool pole", []string{"Burn", "Tarping"}, nil),
	plib.NewBasicItem("ropes", []string{"Burn", "Tarping"}, []string{"Firefly"}),
	plib.NewBasicItem("bungies", []string{"Burn", "Tarping"}, nil),
	plib.NewBasicItem("gaff tape", []string{"Burn", "Tarping"}, nil),
	plib.NewBasicItem("sharpie", []string{"Burn", "Tarping"}, nil),
	plib.NewBasicItem("shims", []string{"Burn", "Tarping"}, nil),
	plib.NewBasicItem("2x battery string lights", []string{"Burn", "Con"}, nil),
	plib.NewBasicItem("camp chair", []string{"Camping"}, nil),
	plib.NewBasicItem("small camp table", []string{"Camping"}, nil),
	plib.NewBasicItem("4x8 camp table", []string{"Burn"}, []string{"Firefly"}),
	plib.NewBasicItem("camp stove", []string{"GrumpCamping"}, []string{"NoFire"}),
	plib.NewBasicItem("propane", []string{"GrumpCamping"}, []string{"NoFire"}),
	plib.NewBasicItem("saucepan", []string{"GrumpCamping"}, nil),
	plib.NewBasicItem("frying pan", []string{"GrumpCamping"}, nil),
	plib.NewTemperatureItem("1 box hot hands", 0, 45, []string{"Camping"}, nil),
	plib.NewBasicItem("camp plateware", []string{"Camping"}, nil),
	plib.NewBasicItem("A CUP", []string{"Burn", "Camping"}, nil),
	plib.NewBasicItem("multitool", []string{"Camping", "Hiking"}, nil),
	plib.NewBasicItem("lighter/matches", []string{"Camping"}, []string{"NoFire", "Firefly"}),
	plib.NewBasicItem("bug spray", []string{"Camping", "Hiking"}, nil),
	plib.NewConsumableItem("paper towels", 0.25, "rolls", []string{"Camping"}, nil),
	plib.NewConsumableItem("trash bag", 0.25, plib.NoUnits, []string{"Camping"}, nil),
	plib.NewConsumableItem("recycle bag", 0.125, plib.NoUnits, []string{"Camping", "Tiny House"}, nil),
	plib.NewBasicItem("camelbak", []string{"Burn", "Camping", "Hiking"}, nil),
	plib.NewBasicItem("extra camelbak bite valves", []string{"Burn", "Camping"}, nil),
	plib.NewBasicItem("cart", []string{"Burn"}, nil),
	plib.NewBasicItem("towel", []string{"Camping", "Tiny House"}, nil),
	plib.NewBasicItem("2x configured radios", []string{"Burn"}, []string{"Firefly"}),
	plib.NewBasicItem("pencils", []string{"Burn"}, nil),
	plib.NewBasicItem("gear ties", []string{"Burn", "Tarping"}, nil),
	plib.NewBasicItem("first aid kit", []string{"Hiking"}, nil),
	plib.NewBasicItem("sheet set", []string{"Tiny House"}, nil),
	plib.NewBasicItem("hammock", []string{"Burn", "Camping"}, nil),
}

func init() {
	plib.RegisterItems("Camping", campStuff)
}
