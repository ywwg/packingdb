package items

import (
	plib "github.com/ywwg/packingdb/pkg/packinglib"
)

var food = []*plib.Item{
	plib.NewItem("booze", []string{"Burn", "BYOB"}, nil).Consumable(0.25),
	plib.NewItem("good beer", []string{"Burn", "BYOB"}, nil).Units("sixpacks").Consumable(0.25),
	plib.NewItem("Drinking water", []string{"Burn"}, nil).Units("gallons").Consumable(0.5),
	// Burns use a camelbak
	plib.NewItem("zojirushi bottle", nil, []string{"Burn"}),
	plib.NewItem("cooking water", []string{"Burn"}, nil).Units("gallons").Consumable(0.25),
	plib.NewItem("cooler", []string{"Camping"}, nil),
	plib.NewItem("spatula", []string{"Camping"}, nil),
	plib.NewItem("tongs", []string{"Camping"}, nil),
	plib.NewItem("wood spoon", []string{"Camping"}, nil),
	plib.NewItem("tasty bites", []string{"Burn"}, nil).Consumable(0.75),
	plib.NewItem("plates", []string{"Con"}, nil),
	plib.NewItem("plastic cutlery", []string{"Con"}, nil),
	plib.NewItem("energy bars", []string{"Camping", "Con", "Cycling", "Hiking"}, nil).Consumable(1.5),
	plib.NewItem("nuun", []string{"Camping", "Cycling", "Hiking"}, nil).Units("tubes").Consumable(.333),
	plib.NewItem("eggs", []string{"Burn"}, nil).Custom(func(_ float64, nights int, _ plib.PropertySet) float64 {
		// 2 Eggs for each morning that I'm there.
		return float64((nights - 1) * 2)
	}),
	plib.NewItem("junk food?", []string{"Camping"}, nil),
	plib.NewItem("salt and pepper", []string{"Camping"}, nil),
	plib.NewItem("box soup", []string{"Camping"}, []string{"NoFire"}),
	plib.NewItem("hot sauce", []string{"Camping"}, nil),
	plib.NewItem("frozen grillables", []string{"Camping"}, []string{"NoFire"}).Units("servings").Consumable(0.75),
	plib.NewItem("buns", []string{"Camping"}, []string{"NoFire"}).Units("servings").Consumable(0.75),

	// Tiny House food
	plib.NewItem("breakfasts", []string{"Tiny House"}, nil).Consumable(1.0),
	plib.NewItem("lunches", []string{"Tiny House"}, nil).Consumable(1.0),
	plib.NewItem("dinners", []string{"Tiny House"}, nil).Consumable(1.0),
	plib.NewItem("all non-freezable ingredients", []string{"Tiny House"}, nil),

	plib.NewItem("green tea", []string{"Lodging", "Con"}, nil).Units("bags").Consumable(1.0).Max(6.0),
	plib.NewItem("no-think food for drive", []string{"Tiny House", "CarRide"}, nil),
	plib.NewItem("liquid for drive", []string{"Tiny House", "CarRide"}, nil),
}
