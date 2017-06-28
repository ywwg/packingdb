package items

import (
	plib "github.com/ywwg/packinglib"
)

var burnDinner = []plib.Item{
	plib.NewConsumableItem("pulled pork", 8.5, "pounds", []string{"Burn"}, nil),
	plib.NewConsumableItem("vegan chili", 3.5, "gallons", []string{"Burn"}, nil),
	plib.NewConsumableItem("pulled jackfruit", 50, "ounces", []string{"Burn"}, nil),
	plib.NewConsumableItem("shredded cheese", 2, "pounds", []string{"Burn"}, nil),
	plib.NewConsumableItem("sour cream", 24, "ounces", []string{"Burn"}, nil),
	plib.NewBasicItem("sarina's cooler", []string{"Camping"}, nil),
}

func init() {
	plib.RegisterItems("Burn Night Dinner", burnDinner)
}
