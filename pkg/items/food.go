package items

import (
	plib "github.com/ywwg/packingdb/pkg/packinglib"
)

var food = []*plib.Item{
	plib.NewItem("booze", []string{"Burn", "BYOB"}, nil).Consumable(0.25, plib.NoUnits),
	plib.NewItem("good beer", []string{"Burn", "BYOB"}, nil).Consumable(0.25, "sixpacks"),
	plib.NewItem("Drinking water", []string{"Burn"}, nil).Consumable(0.5, "gallons"),
	// Burns use a camelbak
	plib.NewItem("zojirushi bottle", nil, []string{"Burn"}),
	plib.NewItem("cooking water", []string{"Burn"}, nil).Consumable(0.25, "gallons"),
	plib.NewItem("cooler", []string{"Camping"}, nil),
	plib.NewItem("spatula", []string{"Camping"}, nil),
	plib.NewItem("tongs", []string{"Camping"}, nil),
	plib.NewItem("wood spoon", []string{"Camping"}, nil),
	plib.NewItem("tasty bites", []string{"Burn"}, nil).Consumable(0.75, plib.NoUnits),
	plib.NewItem("plates", []string{"Con"}, nil),
	plib.NewItem("plastic cutlery", []string{"Con"}, nil),
	plib.NewItem("energy bars", []string{"Camping", "Con", "Cycling", "Hiking"}, nil).Consumable(1.5, plib.NoUnits),
	plib.NewItem("nuun", []string{"Camping", "Cycling", "Hiking"}, nil).Consumable(.333, "tubes"),
	plib.NewItem("eggs", []string{"Burn"}, nil).Custom(func(_ float64, nights int, _ plib.PropertySet) float64 {
		// 2 Eggs for each morning that I'm there.
		return float64((nights - 1) * 2)
	}, plib.NoUnits),
	plib.NewItem("junk food?", []string{"Camping"}, nil),
	plib.NewItem("salt and pepper", []string{"Camping"}, nil),
	plib.NewItem("box soup", []string{"Camping"}, []string{"NoFire"}),
	plib.NewItem("hot sauce", []string{"Camping"}, nil),
	plib.NewItem("frozen grillables", []string{"Camping"}, []string{"NoFire"}).Consumable(0.75, "servings"),
	plib.NewItem("buns", []string{"Camping"}, []string{"NoFire"}).Consumable(0.75, "servings"),

	// Tiny House food
	plib.NewItem("breakfasts", []string{"Tiny House"}, nil).Consumable(1.0, plib.NoUnits),
	plib.NewItem("lunches", []string{"Tiny House"}, nil).Consumable(1.0, plib.NoUnits),
	plib.NewItem("dinners", []string{"Tiny House"}, nil).Consumable(1.0, plib.NoUnits),
	plib.NewItem("all non-freezable ingredients", []string{"Tiny House"}, nil),

	plib.NewItem("green tea", []string{"Lodging", "Con"}, nil).Consumable(1.0, "bags").Max(6.0),
	plib.NewItem("no-think food for drive", []string{"Tiny House", "CarRide"}, nil),
	plib.NewItem("liquid for drive", []string{"Tiny House", "CarRide"}, nil),
}

func init() {
	plib.RegisterItems("Food", food)
}
