package items

import (
	plib "github.com/ywwg/packingdb/pkg/packinglib"
)

var food = []*plib.Item{
	plib.NewConsumableItem("booze", 0.25, plib.NoUnits, []string{"Burn", "BYOB"}, nil),
	plib.NewConsumableItem("good beer", 0.25, "sixpacks", []string{"Burn", "BYOB"}, nil),
	plib.NewConsumableItem("Drinking water", 0.5, "gallons", []string{"Burn"}, nil),
	// Burns use a camelbak
	plib.NewItem("zojirushi bottle", nil, []string{"Burn"}),
	plib.NewConsumableItem("cooking water", 0.25, "gallons", []string{"Burn"}, nil),
	plib.NewItem("cooler", []string{"Camping"}, nil),
	plib.NewItem("spatula", []string{"Camping"}, nil),
	plib.NewItem("tongs", []string{"Camping"}, nil),
	plib.NewItem("wood spoon", []string{"Camping"}, nil),
	plib.NewConsumableItem("tasty bites", 0.75, plib.NoUnits, []string{"Burn"}, nil),
	plib.NewItem("plates", []string{"Con"}, nil),
	plib.NewItem("plastic cutlery", []string{"Con"}, nil),
	plib.NewConsumableItem("energy bars", 1.5, plib.NoUnits, []string{"Camping", "Con", "Cycling", "Hiking"}, nil),
	plib.NewConsumableItem("nuun", .333, "tubes", []string{"Camping", "Cycling", "Hiking"}, nil),
	plib.NewCustomConsumableItem("eggs", func(nights int, _ plib.PropertySet) float64 {
		// 2 Eggs for each morning that I'm there.
		return float64((nights - 1) * 2)
	}, plib.NoUnits, []string{"Burn"}, nil),
	plib.NewItem("junk food?", []string{"Camping"}, nil),
	plib.NewItem("salt and pepper", []string{"Camping"}, nil),
	plib.NewItem("box soup", []string{"Camping"}, []string{"NoFire"}),
	plib.NewItem("hot sauce", []string{"Camping"}, nil),
	plib.NewConsumableItem("frozen grillables", 0.75, "servings", []string{"Camping"}, []string{"NoFire"}),
	plib.NewConsumableItem("buns", 0.75, "servings", []string{"Camping"}, []string{"NoFire"}),

	// Tiny House food
	plib.NewConsumableItem("breakfasts", 1.0, plib.NoUnits, []string{"Tiny House"}, nil),
	plib.NewConsumableItem("lunches", 1.0, plib.NoUnits, []string{"Tiny House"}, nil),
	plib.NewConsumableItem("dinners", 1.0, plib.NoUnits, []string{"Tiny House"}, nil),
	plib.NewItem("all non-freezable ingredients", []string{"Tiny House"}, nil),

	plib.NewConsumableMaxItem("green tea", 1.0, 6.0, "bags", []string{"Lodging", "Con"}, nil),
	plib.NewItem("no-think food for drive", []string{"Tiny House", "CarRide"}, nil),
	plib.NewItem("liquid for drive", []string{"Tiny House", "CarRide"}, nil),
}

func init() {
	plib.RegisterItems("Food", food)
}
