package items

import (
	plib "github.com/ywwg/packinglib"
)

var food = []plib.Item{
	plib.NewConsumableItem("booze", 0.25, plib.NoUnits, []string{"Burn"}, nil),
	plib.NewConsumableItem("good beer", 0.25, "sixpacks", []string{"Burn"}, nil),
	plib.NewConsumableItem("drinking water", 0.5, "gallons", []string{"Burn"}, nil),
	plib.NewConsumableItem("cooking water", 0.25, "gallons", []string{"Burn"}, nil),
	plib.NewBasicItem("cooler", []string{"Camping"}, nil),
	plib.NewConsumableItem("tasty bites", 0.75, plib.NoUnits, []string{"Burn"}, nil),
	plib.NewConsumableItem("clif bars", 1.5, plib.NoUnits, []string{"Camping"}, nil),
	plib.NewConsumableItem("nuun", .333, "tubes", []string{"Camping"}, nil),
	plib.NewCustomConsumableItem("eggs", func(nights int) float64 {
		// 2 Eggs for each morning that I'm there.
		return float64((nights - 1) * 2)
	}, plib.NoUnits, []string{"Burn"}, nil),
	plib.NewBasicItem("mustard", []string{"Camping"}, nil),
	plib.NewBasicItem("junk food?", []string{"Camping"}, nil),
	plib.NewBasicItem("lamb biryani", []string{"Firefly2017"}, nil),
	plib.NewBasicItem("salt and pepper", []string{"Camping"}, nil),
	plib.NewBasicItem("box soup", []string{"Camping"}, nil),
	plib.NewBasicItem("hot sauce", []string{"Camping"}, nil),
	plib.NewConsumableItem("frozen grillables", 0.75, "servings", []string{"Camping"}, nil),
	plib.NewConsumableItem("buns", 0.75, "servings", []string{"Camping"}, nil),

	// Tiny House food
	plib.NewConsumableItem("breakfasts", 1.0, plib.NoUnits, []string{"Tiny House"}, nil),
	plib.NewConsumableItem("lunches", 1.0, plib.NoUnits, []string{"Tiny House"}, nil),
	plib.NewConsumableItem("dinners", 1.0, plib.NoUnits, []string{"Tiny House"}, nil),
	plib.NewBasicItem("all non-freezable ingredients", []string{"Tiny House"}, nil),
}

func init() {
	plib.RegisterItems("Food", food)
}
