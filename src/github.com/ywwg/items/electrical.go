package items

import (
	plib "github.com/ywwg/packinglib"
)

var electrical = []plib.Item{
	plib.NewBasicItem("Generator", []string{"Burn"}, nil),
	plib.NewCustomConsumableItem("Gas", func(days int) float64 {
		return float64(days - 1)
	}, "gallons", []string{"Burn"}, nil),
	plib.NewBasicItem("Kill-A-Watt", []string{"Burn"}, nil),
	plib.NewBasicItem("Bluetooth speaker and cord", nil, nil),
	plib.NewBasicItem("long extension cord", []string{"Burn"}, nil),
	plib.NewBasicItem("Cell Phone Charger", nil, nil),
	plib.NewBasicItem("Batteries for headlamp", nil, nil),
	plib.NewBasicItem("Sound meter", []string{"Burn"}, nil),
}

func init() {
	plib.RegisterItems("Electrical", electrical)
}
