package items

import (
	plib "github.com/ywwg/packinglib"
)

var clothing = []plib.Item{
	plib.NewBasicItem("boots", []string{"Dirt"}, nil),
	plib.NewConsumableItem("underwear", 1.0, "pair", nil, nil),
	plib.NewConsumableItem("crew socks", 1.0, "pair", nil, []string{"Camping"}),
	plib.NewConsumableItem("smartwool socks", 1.0, "pair", []string{"Camping"}, nil),
	plib.NewConsumableTemperatureItem("long pj socks", 0.25, "pair", 0, 55, nil, nil),
	plib.NewConsumableTemperatureItem("short pj socks", 0.25, "pair", 60, 120, nil, nil),
	plib.NewConsumableTemperatureItem("short pjs", 0.2, "set", 75, 120, nil, nil),
	plib.NewConsumableTemperatureItem("long pjs", 0.2, "set", 0, 74, nil, nil),
	plib.NewConsumableItem("jeans", 0.2, "pair", nil, nil),
	plib.NewConsumableItem("convertapants", 0.5, "pair", []string{"Dirt"}, nil),
	plib.NewConsumableTemperatureItem("long underwear", 0.25, "pair", 0, 45, []string{"Camping"}, nil),
	plib.NewConsumableTemperatureItem("tshirts", 0.75, plib.NoUnits, 65, 120, nil, nil),
	plib.NewConsumableTemperatureItem("non-cotton tops", 0.25, plib.NoUnits, 45, 65, []string{"Camping"}, nil),
	plib.NewConsumableTemperatureItem("underlayer top", 0.5, plib.NoUnits, 0, 50, []string{"Camping"}, nil),
	plib.NewConsumableTemperatureItem("underlayer bottom", 0.5, plib.NoUnits, 0, 50, []string{"Camping"}, nil),
	plib.NewConsumableItem("fun outfits", 0.75, plib.NoUnits, []string{"Partying"}, nil),
	plib.NewBasicItem("kigarumi", []string{"Partying"}, nil),
	plib.NewConsumableTemperatureItem("sweaty shirts", 0.5, plib.NoUnits, 65, 120, []string{"Dirt"}, nil),
	plib.NewTemperatureItem("light jacket", 55, 65, nil, nil),
	plib.NewTemperatureItem("medium jacket", 40, 55, nil, nil),
	plib.NewTemperatureItem("heavy jacket", 0, 40, nil, nil),
	plib.NewBasicItem("rain jacket", []string{"Camping"}, nil),
	plib.NewTemperatureItem("warm hat", 0, 40, nil, nil),
	plib.NewTemperatureItem("gloves", 0, 40, nil, nil),
	plib.NewBasicItem("earplugs", []string{"Loud"}, nil),
	plib.NewBasicItem("eye mask", []string{"Loud"}, nil),
	plib.NewTemperatureItem("tevas", 70, 120, nil, nil),
	plib.NewBasicItem("nice dinner clothes", []string{"Fancy"}, nil),
	plib.NewBasicItem("watch", nil, nil),
}

func init() {
	plib.RegisterItems("Clothing", clothing)
}
