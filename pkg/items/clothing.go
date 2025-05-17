package items

import (
	plib "github.com/ywwg/packingdb/pkg/packinglib"
)

var clothing = []*plib.Item{
	plib.NewItem("sleep earplugs", nil, nil),
	plib.NewItem("eye mask", nil, nil),
	plib.NewItem("face masks", nil, nil),
	plib.NewItem("boots", []string{"Dirt", "Snow", "Tiny House", "Skiing", "Hiking"}, nil),
	plib.NewItem("underwear", nil, nil).Units("pair").Consumable(1.0),
	plib.NewItem("crew socks", nil, []string{"Camping"}).Units("pair").TemperatureRange(0, 65).Consumable(1.0).Max(4.0),
	plib.NewItem("short socks", nil, []string{"Camping"}).Units("pair").TemperatureRange(70, 120).Consumable(1.0).Max(4.0),
	plib.NewItem("sneakers", nil, []string{"Burn"}),
	plib.NewItem("smartwool socks", []string{"Camping", "Tiny House", "Skiing", "Hiking"}, nil).Units("pair").Consumable(1.0),
	plib.NewItem("short pj socks", nil, nil).TemperatureRange(60, 120).Units("pair"),
	plib.NewItem("short pjs", nil, nil).TemperatureRange(70, 120),
	plib.NewItem("long pjs", nil, nil).TemperatureRange(0, 55),
	plib.NewItem("work jeans", []string{"Dirt"}, []string{"Burn"}).Units("pair").Consumable(0.2).Max(2.0),
	plib.NewItem("shorts", nil, []string{"Burn"}).Units("pair").TemperatureRange(70, 120).Consumable(0.3),
	plib.NewItem("jeans", nil, nil).TemperatureRange(0, 70).Units("pair").Consumable(0.1),
	plib.NewItem("pristine shirt for drive home", []string{"Burn"}, nil),
	plib.NewItem("convertapants", []string{"Dirt", "Burn", "Hiking"}, nil).Units("pair").Consumable(0.5).Max(2.0),
	plib.NewItem("sweater", nil, nil).TemperatureRange(0, 55).Consumable(0.5),
	plib.NewItem("undershirt", nil, nil).TemperatureRange(0, 50).Consumable(0.5),
	plib.NewItem("long underwear", []string{"Camping", "Tiny House"}, nil).Units("pair").TemperatureRange(0, 45).Consumable(0.25),
	plib.NewItem("thick slipper socks", nil, nil).TemperatureRange(0, 50).Consumable(0.25),
	plib.NewItem("tshirts", nil, nil).TemperatureRange(40, 120).Consumable(0.75).Max(4.0),
	plib.NewItem("non-cotton tops", []string{"Camping", "Hiking"}, nil).TemperatureRange(45, 120).Consumable(0.25),
	plib.NewItem("underlayer top", []string{"Camping", "Skiing", "Hiking"}, nil).TemperatureRange(0, 50).Consumable(0.5),
	plib.NewItem("underlayer bottom", []string{"Camping", "Skiing", "Hiking"}, nil).TemperatureRange(0, 50).Consumable(0.5),
	plib.NewItem("fun outfits", []string{"Partying"}, nil).Consumable(0.75).Max(3.0),
	plib.NewItem("shirt for flight", []string{"Flight"}, nil),
	plib.NewItem("sweaty shirts", []string{"Dirt", "Handy"}, nil).TemperatureRange(65, 120).Consumable(0.5).Max(3.0),
	plib.NewItem("light jacket", nil, nil).TemperatureRange(51, 60),
	plib.NewItem("medium jacket", nil, nil).TemperatureRange(41, 50),
	plib.NewItem("heavy jacket", nil, nil).TemperatureRange(0, 40),
	plib.NewItem("rain jacket", nil, nil),
	plib.NewItem("snow pants", []string{"Snow", "Skiing"}, nil),
	plib.NewItem("umbrella", nil, nil),
	plib.NewItem("warm hat", nil, nil).TemperatureRange(0, 50),
	plib.NewItem("scarf", nil, nil).TemperatureRange(0, 40),
	plib.NewItem("gloves", nil, nil).TemperatureRange(0, 50),
	plib.NewItem("tevas", nil, []string{"International"}).TemperatureRange(70, 120),
	plib.NewItem("nice dinner clothes", []string{"DiningOut"}, nil),
	plib.NewItem("watch", nil, nil),
	plib.NewItem("purse", nil, []string{"Burn", "Camping"}),

	plib.NewItem("three piece suit", []string{"UltraFormal"}, nil),
	plib.NewItem("thin belt", []string{"UltraFormal"}, nil),
	plib.NewItem("suit shirts", []string{"UltraFormal"}, nil).Consumable(0.5),
	plib.NewItem("ties", []string{"UltraFormal"}, nil).Consumable(0.5),
}

var flightSupplies = []*plib.Item{
	plib.NewItem("turtl pillow", []string{"Flight", "International", "PaidTravel"}, nil),
	plib.NewItem("compression socks in carryon", []string{"Flight", "International"}, nil),
}

var conSupplies = []*plib.Item{
	plib.NewItem("fanny pack", []string{"Con", "Partying"}, nil),
}
