package items

import (
	plib "github.com/ywwg/packinglib"
)

var campStuff = []plib.Item{
	plib.NewBasicItem("tent", []string{"Camping"}, nil),
	plib.NewBasicItem("sleeping bag", []string{"Camping"}, nil),
	plib.NewBasicItem("sleeping pad", []string{"Camping"}, nil),
	plib.NewBasicItem("tent light", []string{"Camping"}, nil),
	plib.NewBasicItem("headlamp", []string{"Dark"}, nil),
	plib.NewBasicItem("towel for butt driving home", []string{"Burn"}, nil),
	plib.NewBasicItem("tarps", []string{"Burn"}, nil),
	plib.NewBasicItem("ropes", []string{"Burn"}, nil),
	plib.NewBasicItem("lighter/matches", []string{"Camping"}, nil),
	plib.NewBasicItem("bug spray", []string{"Camping"}, nil),
	plib.NewBasicItem("dirty clothes bag", []string{"Camping"}, nil),
	plib.NewConsumableItem("trash bag", 0.25, plib.NoUnits, []string{"Camping"}, nil),
	plib.NewConsumableItem("recycle bag", 0.125, plib.NoUnits, []string{"Camping"}, nil),
	plib.NewConsumableItem("drinking water", 0.5, "gallons", []string{"Burn"}, nil),
	plib.NewConsumableItem("cooking water", 0.25, "gallons", []string{"Burn"}, nil),
	plib.NewBasicItem("camelbak", []string{"Burn", "Camping"}, nil),
	plib.NewBasicItem("cart", []string{"Burn"}, nil),
	plib.NewBasicItem("TASK: permetherin", []string{"Camping"}, nil),
}

func init() {
	plib.RegisterItems("Camping Stuff", campStuff)
}
