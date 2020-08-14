package items

import (
	plib "github.com/ywwg/packinglib"
)

var tools = []plib.Item{
	plib.NewBasicItem("Tool Bag", []string{"Handy"}, nil),
	plib.NewBasicItem("Tiny House ALC", []string{"Tiny House"}, nil),
	plib.NewConsumableItem("Gas for Generator", 1.0, "gallons", []string{"Tiny House"}, nil),
	plib.NewBasicItem("Hand vac", []string{"Tiny House"}, nil),
	plib.NewBasicItem("5 gallon water jug", []string{"Tiny House"}, nil),
	plib.NewBasicItem("UDC", []string{"Tiny House"}, nil),
}

func init() {
	plib.RegisterItems("Tools", tools)
}
