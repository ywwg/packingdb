package items

import (
	plib "github.com/ywwg/packingdb/pkg/packinglib"
)

var tools = []*plib.Item{
	plib.NewItem("Tool Bag", []string{"Handy"}, nil),
	plib.NewItem("Tiny House ALC", []string{"Tiny House"}, nil),
	plib.NewItem("Gas for Generator", []string{"Tiny House"}, nil).Consumable(1.0, "gallons"),
	plib.NewItem("Hand vac", []string{"Tiny House"}, nil),
	plib.NewItem("5 gallon water jug", []string{"Tiny House"}, nil),
	plib.NewItem("UDC", []string{"Tiny House"}, nil),
}

func init() {
	plib.RegisterItems("Tools", tools)
}
