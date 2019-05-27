package items

import (
	plib "github.com/ywwg/packinglib"
)

var electrical = []plib.Item{
	plib.NewBasicItem("Generator", []string{"Burn", "Tiny House"}, nil),
	plib.NewBasicItem("Fire extinguisher for genny", []string{"Burn"}, nil),
	plib.NewConsumableItem("Gas", 1.0, "gallons", []string{"Burn"}, nil),
	plib.NewBasicItem("Kill-A-Watt", []string{"Burn"}, nil),
	plib.NewBasicItem("1/8 to 1/8", []string{"Performing"}, []string{"Tiny House"}),
	plib.NewBasicItem("long extension cord", []string{"Burn"}, nil),
	plib.NewBasicItem("cell phone fast charger", nil, []string{"Tiny House"}),
	plib.NewBasicItem("battery pack and microsd charger", nil, []string{"Tiny House"}),
	plib.NewBasicItem("Batteries for headlamp", []string{"Burn"}, nil),
	plib.NewBasicItem("DSLR and Charger", []string{"Camping", "BigTrip"}, []string{"Insecure"}),
	plib.NewBasicItem("Plug Adapters", []string{"International"}, nil),
}

func init() {
	plib.RegisterItems("Electrical", electrical)
}
