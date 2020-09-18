package items

import (
	plib "github.com/ywwg/packinglib"
)

var dog = []plib.Item{
	plib.NewBasicItem("6ft Leash", []string{"Dog"}, nil),
	plib.NewBasicItem("Dog harness", []string{"Dog"}, nil),
	plib.NewBasicItem("Long tie-down rope", []string{"Dog"}, nil),
	plib.NewBasicItem("Dog toys", []string{"Dog"}, nil),
	plib.NewBasicItem("Collar light", []string{"Dog"}, nil),
	plib.NewBasicItem("Dog towel", []string{"Dog"}, nil),
	plib.NewBasicItem("Dog water bowl", []string{"Dog"}, nil),
	plib.NewBasicItem("Dog food bowl", []string{"Dog"}, nil),
	plib.NewConsumableItem("Dog food", 2.0, "servings", []string{"Dog"}, nil),
	plib.NewConsumableItem("Dog treats", 1.0, "days worth", []string{"Dog"}, nil),
	plib.NewConsumableItem("rawhide", 1.0, "sticks", []string{"Dog"}, nil),
	plib.NewBasicItem("Dog food bag in car", []string{"Dog"}, nil),
}

func init() {
	plib.RegisterItems("Dog", dog)
}
