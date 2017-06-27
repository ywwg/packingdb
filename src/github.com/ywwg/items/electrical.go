package items

import (
	plib "github.com/ywwg/packinglib"
)

var electrical = []plib.Item{
	plib.NewBasicItem("Generator", []string{"Burn"}, nil),
	plib.NewConsumableItem("Gas", 1.0, "gallon", []string{"Burn"}, nil),
	plib.NewBasicItem("Extension Cords", []string{"Burn"}, nil),
	plib.NewBasicItem("Cell Phone Charger", nil, nil),
}

func init() {
	plib.RegisterItems("Electrical", electrical)
}
