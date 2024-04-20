package items

import (
	plib "github.com/ywwg/packingdb/pkg/packinglib"
)

var electrical = []*plib.Item{
	plib.NewItem("Generator", []string{"Burn", "Tiny House"}, nil),
	plib.NewItem("Fire extinguisher for genny", []string{"Burn"}, nil),
	plib.NewConsumableItem("Gas", 1.0, "gallons", []string{"Burn"}, nil),
	plib.NewItem("Kill-A-Watt", []string{"Burn"}, nil),
	plib.NewItem("long extension cord", []string{"Burn"}, nil),
	plib.NewItem("cell phone fast charger", nil, []string{"Tiny House"}),
	plib.NewItem("battery pack and microsd charger", nil, []string{"Tiny House"}),
	plib.NewItem("Batteries for headlamp", []string{"Burn"}, nil),
	plib.NewItem("DSLR and Charger", []string{"Camping", "BigTrip"}, []string{"Insecure"}),
	plib.NewItem("Plug Adapters", []string{"International"}, nil),
	plib.NewItem("Disco laser", []string{"Con"}, nil),
}

func init() {
	plib.RegisterItems("Electrical", electrical)
}
