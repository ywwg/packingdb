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
	plib.NewConsumableItem("clif bars", 1.5, plib.NoUnits, []string{"Burn"}, nil),
	plib.NewCustomConsumableItem("eggs", func(days int) float64 {
		// 2 Eggs for each morning that I'm there.
		return float64((days - 1) * 2)
	}, plib.NoUnits, []string{"Burn"}, nil),
	plib.NewBasicItem("mustard", []string{"Camping"}, nil),
	plib.NewBasicItem("salt and pepper", []string{"Camping"}, nil),
	plib.NewBasicItem("box soup", []string{"Camping"}, nil),
	plib.NewBasicItem("hot sauce", []string{"Camping"}, nil),
	plib.NewConsumableItem("grillables", 0.75, "servings", []string{"Camping"}, nil),
}

func init() {
	plib.RegisterItems("Food", food)
}
