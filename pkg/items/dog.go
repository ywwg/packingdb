package items

import (
	plib "github.com/ywwg/packingdb/pkg/packinglib"
)

var dog = []*plib.Item{
	plib.NewItem("6ft Leash", []string{"Dog"}, nil),
	plib.NewItem("Dog harness", []string{"Dog"}, nil),
	plib.NewItem("Long tie-down rope", []string{"Dog"}, nil),
	plib.NewItem("Dog toys", []string{"Dog"}, nil),
	plib.NewItem("Collar light", []string{"Dog"}, nil),
	plib.NewItem("Dog towel", []string{"Dog"}, nil),
	plib.NewItem("Dog water bowl", []string{"Dog"}, nil),
	plib.NewItem("Dog food bowl", []string{"Dog"}, nil),
	plib.NewItem("Dog food dry", []string{"Dog"}, nil).Consumable(2.0, "servings"),
	plib.NewItem("Dog food wet", []string{"Dog"}, nil).Consumable(2.0, "servings"),
	plib.NewItem("Dog treats", []string{"Dog"}, nil).Consumable(1.0, "days worth"),
	plib.NewItem("rawhide", []string{"Dog"}, nil).Consumable(1.0, "sticks"),
	plib.NewItem("Dog food bag in car", []string{"Dog"}, nil),
}

func init() {
	plib.RegisterItems("Dog", dog)
}
