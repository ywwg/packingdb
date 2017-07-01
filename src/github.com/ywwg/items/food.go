package items

import (
	plib "github.com/ywwg/packinglib"
)

var food = []plib.Item{
	plib.NewConsumableItem("booze", 0.5, plib.NoUnits, []string{"Burn"}, nil),
	plib.NewConsumableItem("good beer", 0.33, "sixpacks", []string{"Burn"}, nil),
	plib.NewConsumableItem("drinking water", 0.5, "gallons", []string{"Burn"}, nil),
	plib.NewConsumableItem("cooking water", 0.25, "gallons", []string{"Burn"}, nil),
	plib.NewBasicItem("cooler", []string{"Camping"}, nil),
	plib.NewConsumableItem("tasty bites", 1.0, plib.NoUnits, []string{"Burn"}, nil),
	plib.NewConsumableItem("clif bars", 1.5, plib.NoUnits, []string{"Burn"}, nil),
	plib.NewCustomConsumableItem("eggs", func(days int) float64 {
		// 2 Eggs for each morning that I'm there.
		return float64((days - 1) * 2)
	}, plib.NoUnits, []string{"Burn"}, nil),
}

func init() {
	plib.RegisterItems("Food", food)
}
