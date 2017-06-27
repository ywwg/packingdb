package packinglib

//var Clothing []Item

var Clothing = []Item{
	NewBasicItem("boots", PropertySet{Dirt: true}),
	NewConsumableItem("underwear", 1.0, "pair", nil),
}

//	{"crew socks", 1.0, "pair", 0, 120, nil},
//	{"short pjs", 0.2, "set", 75, 120, nil},
//	{"long pjs", 0.2, "set", 0, 74, nil},
//	{"jeans", 0.2, "pair", 0, 120, nil},
//	{"convertapants", 0.5, "pair", 0, 120, PropertySet{Dirt: true}},
//	{"long underwear", 0.25, "pair", 0, 45, PropertySet{Camping: true}},
//	{"tshirts", 0.75, NoUnits, 65, 120, nil},
//	{"sweaty shirts", 0.5, NoUnits, 65, 120, PropertySet{Dirt: true}},
//	{"light jacket", 0.0, NoUnits, 55, 65, nil},
//	{"medium jacket", 0.0, NoUnits, 40, 55, nil},
//	{"heavy jacket", 0.0, NoUnits, 0, 40, nil},
//	{"rain jacket", 0.0, NoUnits, 0, 120, PropertySet{Camping: true}},
//	{"warm hat", 0.0, NoUnits, 0, 40, nil},
//	{"gloves", 0.0, NoUnits, 0, 40, nil},
//	{"earplugs", 0.0, NoUnits, 0, 120, PropertySet{Loud: true}},
//	{"face mask", 0.0, NoUnits, 0, 120, PropertySet{Loud: true}},
//	{"tevas", 0.0, NoUnits, 70, 120, nil},
//	{"boots", 0.0, NoUnits, 0, 120, PropertySet{Dirt: true}},
//	{"thick socks", 0.5, "pair", 0, 50, PropertySet{Dirt: true}},
//	{"nice dinner clothes", 0.0, NoUnits, 0, 120, PropertySet{Fancy: true}},
//	//
//	{"sunscreen", 0.0, NoUnits, 65, 120, nil},
//}

func init() {
	//	Clothing = append(Clothing, NewBasicItem("boots", PropertySet{Dirt: true}))
	//	 Clothing = []Item{
	//	{"underwear", 1.0, "pair", 0, 120, nil},
	//	{"crew socks", 1.0, "pair", 0, 120, nil},
	//	{"short pjs", 0.2, "set", 75, 120, nil},
	//	{"long pjs", 0.2, "set", 0, 74, nil},
	//	{"jeans", 0.2, "pair", 0, 120, nil},
	//	{"convertapants", 0.5, "pair", 0, 120, PropertySet{Dirt: true}},
	//	{"long underwear", 0.25, "pair", 0, 45, PropertySet{Camping: true}},
	//	{"tshirts", 0.75, NoUnits, 65, 120, nil},
	//	{"sweaty shirts", 0.5, NoUnits, 65, 120, PropertySet{Dirt: true}},
	//	{"light jacket", 0.0, NoUnits, 55, 65, nil},
	//	{"medium jacket", 0.0, NoUnits, 40, 55, nil},
	//	{"heavy jacket", 0.0, NoUnits, 0, 40, nil},
	//	{"rain jacket", 0.0, NoUnits, 0, 120, PropertySet{Camping: true}},
	//	{"warm hat", 0.0, NoUnits, 0, 40, nil},
	//	{"gloves", 0.0, NoUnits, 0, 40, nil},
	//	{"earplugs", 0.0, NoUnits, 0, 120, PropertySet{Loud: true}},
	//	{"face mask", 0.0, NoUnits, 0, 120, PropertySet{Loud: true}},
	//	{"tevas", 0.0, NoUnits, 70, 120, nil},
	//	{"boots", 0.0, NoUnits, 0, 120, },
	//	{"thick socks", 0.5, "pair", 0, 50, PropertySet{Dirt: true}},
	//	{"nice dinner clothes", 0.0, NoUnits, 0, 120, PropertySet{Fancy: true}},
	//	//
	//	{"sunscreen", 0.0, NoUnits, 65, 120, nil},
	//}
}
