package items

import (
	plib "github.com/ywwg/packinglib"
)

var campStuff = []plib.Item{
	plib.NewBasicItem("windsheild sun shade", []string{"Camping"}, nil),
	plib.NewBasicItem("tent", []string{"Camping"}, nil),
	plib.NewBasicItem("sleeping bag", []string{"Camping"}, nil),
	plib.NewBasicItem("sleeping pad", []string{"Camping"}, nil),
	plib.NewBasicItem("tent light", []string{"Camping"}, nil),
	plib.NewBasicItem("headlamp", []string{"Dark"}, nil),
	plib.NewBasicItem("towel for butt driving home", []string{"Burn"}, nil),
	plib.NewBasicItem("tarps", []string{"Burn"}, nil),
	plib.NewBasicItem("ropes", []string{"Burn"}, nil),
	plib.NewBasicItem("bungies", []string{"Burn"}, nil),
	plib.NewBasicItem("gaff tape", []string{"Burn"}, nil),
	plib.NewBasicItem("camp chair", []string{"Camping"}, nil),
	plib.NewBasicItem("camp table", []string{"Camping"}, nil),
	plib.NewBasicItem("camp stove", []string{"Camping"}, nil),
	plib.NewBasicItem("1 box hot hands?", []string{"Camping"}, nil),
	plib.NewBasicItem("camp plateware", []string{"Camping"}, nil),
	plib.NewBasicItem("A CUP", []string{"Burn"}, nil),
	plib.NewBasicItem("multitool", []string{"Camping"}, nil),
	plib.NewBasicItem("lighter/matches", []string{"Camping"}, nil),
	plib.NewBasicItem("bug spray", []string{"Camping"}, nil),
	plib.NewBasicItem("dirty clothes bag", []string{"Camping"}, nil),
	plib.NewConsumableItem("trash bag", 0.25, plib.NoUnits, []string{"Camping"}, nil),
	plib.NewConsumableItem("recycle bag", 0.125, plib.NoUnits, []string{"Camping"}, nil),
	plib.NewBasicItem("camelbak", []string{"Burn", "Camping"}, nil),
	plib.NewBasicItem("cart", []string{"Burn"}, nil),
	plib.NewBasicItem("TASK: permetherin", []string{"Camping"}, nil),
	plib.NewBasicItem("camp towel", []string{"Camping"}, nil),
	plib.NewBasicItem("2x configured radios", []string{"Burn"}, nil),
}

func init() {
	plib.RegisterItems("Camping Stuff", campStuff)
}
