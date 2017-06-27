package packinglib

//var Clothing []Item

var Clothing = []Item{
	NewBasicItem("boots", PropertySet{Dirt: true}),
	NewConsumableItem("underwear", 1.0, "pair", nil),
	NewConsumableItem("crew socks", 1.0, "pair", nil),
	NewConsumableTemperatureItem("short pjs", 0.2, "set", 75, 120, nil),
	NewConsumableTemperatureItem("long pjs", 0.2, "set", 0, 74, nil),
	NewConsumableItem("jeans", 0.2, "pair", nil),
	NewConsumableItem("convertapants", 0.5, "pair", PropertySet{Dirt: true}),
	NewConsumableTemperatureItem("long underwear", 0.25, "pair", 0, 45, PropertySet{Camping: true}),
	NewConsumableTemperatureItem("tshirts", 0.75, NoUnits, 65, 120, nil),
	NewConsumableTemperatureItem("sweaty shirts", 0.5, NoUnits, 65, 120, PropertySet{Dirt: true}),
	NewTemperatureItem("light jacket", 55, 65, nil),
	NewTemperatureItem("medium jacket", 40, 55, nil),
	NewTemperatureItem("heavy jacket", 0, 40, nil),
	NewBasicItem("rain jacket", PropertySet{Camping: true}),
	NewTemperatureItem("warm hat", 0, 40, nil),
	NewTemperatureItem("gloves", 0, 40, nil),
	NewBasicItem("earplugs", PropertySet{Loud: true}),
	NewBasicItem("face mask", PropertySet{Loud: true}),
	NewTemperatureItem("tevas", 70, 120, nil),
	NewBasicItem("boots", PropertySet{Dirt: true}),
	NewConsumableTemperatureItem("thick socks", 0.5, "pair", 0, 150, PropertySet{Dirt: true}),
	NewBasicItem("nice dinner clothes", PropertySet{Fancy: true}),
	NewTemperatureItem("sunscreen", 65, 120, nil),
}
