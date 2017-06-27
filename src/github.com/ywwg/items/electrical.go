package items

import (
	plib "github.com/ywwg/packinglib"
)

var electrical = []plib.Item{
	plib.NewBasicItem("Generator", []string{"Burn"}, nil),
	plib.NewCustomConsumableItem("Gas", func(days int) float64 {
		return float64(days - 1)
	}, "gallons", []string{"Burn"}, nil),
	plib.NewBasicItem("Extension Cords", []string{"Burn"}, nil),
	plib.NewBasicItem("Cell Phone Charger", nil, nil),
	plib.NewBasicItem("Batteries for headlamp", nil, nil),
}

func init() {
	plib.RegisterItems("Electrical", electrical)
}
