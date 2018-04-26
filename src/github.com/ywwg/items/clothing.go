package items

import (
	plib "github.com/ywwg/packinglib"
)

var clothing = []plib.Item{
	plib.NewBasicItem("boots", []string{"Dirt", "Tiny House"}, nil),
	plib.NewConsumableItem("underwear", 1.0, "pair", nil, nil),
	plib.NewConsumableMaxTemperatureItem("crew socks", 1.0, 4.0, "pair", 0, 65, nil, []string{"Camping"}),
	plib.NewConsumableMaxTemperatureItem("short socks", 1.0, 4.0, "pair", 70, 120, nil, []string{"Camping"}),
	plib.NewBasicItem("sneakers", nil, []string{"Burn"}),
	plib.NewBasicItem("black sneakers", []string{"Berlin2017"}, nil),
	plib.NewConsumableItem("smartwool socks", 1.0, "pair", []string{"Camping", "Tiny House"}, nil),
	plib.NewTemperatureItem("long pj socks", 0, 55, nil, nil),
	plib.NewTemperatureItem("short pj socks", 60, 120, nil, nil),
	plib.NewTemperatureItem("short pjs", 70, 120, nil, nil),
	plib.NewTemperatureItem("long pjs", 0, 55, nil, nil),
	plib.NewConsumableMaxItem("work jeans", 0.2, 2.0, "pair", []string{"Dirt"}, []string{"Burn"}),
	plib.NewConsumableTemperatureItem("shorts", 0.3, "pair", 70, 120, nil, []string{"Burn"}),
	plib.NewConsumableTemperatureItem("jeans", 0.1, "pair", 0, 70, nil, nil),
	plib.NewTemperatureItem("casual belt", 0, 70, nil, nil),
	plib.NewTemperatureItem("fabric belt", 70, 120, nil, nil),
	plib.NewBasicItem("pristine shirt for drive home", []string{"Burn"}, nil),
	plib.NewConsumableItem("convertapants", 0.5, "pair", []string{"Dirt"}, nil),
	plib.NewConsumableTemperatureItem("sweater", 0.5, plib.NoUnits, 0, 50, nil, nil),
	plib.NewConsumableTemperatureItem("undershirt", 0.5, plib.NoUnits, 0, 50, nil, nil),
	plib.NewConsumableTemperatureItem("long underwear", 0.25, "pair", 0, 45, []string{"Camping"}, nil),
	plib.NewConsumableMaxTemperatureItem("tshirts", 0.75, 4.0, plib.NoUnits, 40, 120, nil, nil),
	plib.NewTemperatureItem("longsleeves for under tshirt", 40, 60, nil, nil),
	plib.NewConsumableTemperatureItem("non-cotton tops", 0.25, plib.NoUnits, 45, 120, []string{"Camping"}, nil),
	plib.NewConsumableTemperatureItem("underlayer top", 0.5, plib.NoUnits, 0, 50, []string{"Camping"}, nil),
	plib.NewConsumableTemperatureItem("underlayer bottom", 0.5, plib.NoUnits, 0, 50, []string{"Camping"}, nil),
	plib.NewConsumableMaxItem("fun outfits", 0.75, 3.0, plib.NoUnits, []string{"Partying"}, nil),
	plib.NewBasicItem("REALLY fun outfits", []string{"Berlin2017"}, nil),
	plib.NewBasicItem("shirt for flight", []string{"Flight"}, nil),
	plib.NewBasicItem("kigarumi", []string{"Burn"}, []string{"Firefly2017"}),
	plib.NewConsumableTemperatureItem("sweaty shirts", 0.5, plib.NoUnits, 65, 120, []string{"Dirt"}, nil),
	plib.NewTemperatureItem("light jacket", 51, 60, nil, nil),
	plib.NewTemperatureItem("medium jacket", 41, 50, nil, nil),
	plib.NewTemperatureItem("heavy jacket", 0, 40, nil, nil),
	plib.NewBasicItem("rain jacket", nil, nil),
	plib.NewTemperatureItem("warm hat", 0, 40, nil, nil),
	plib.NewTemperatureItem("gloves", 0, 40, nil, nil),
	plib.NewBasicItem("sleep earplugs", nil, nil),
	plib.NewBasicItem("eye mask", nil, nil),
	plib.NewTemperatureItem("tevas", 70, 120, nil, []string{"International"}),
	plib.NewBasicItem("nice dinner clothes", []string{"Fancy"}, nil),
	plib.NewBasicItem("watch", nil, nil),
	plib.NewBasicItem("purse", nil, []string{"Burn", "Camping"}),
}

func init() {
	plib.RegisterProperty("Berlin2017")
	plib.RegisterItems("Clothing", clothing)
}
