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
	plib.NewBasicItem("Extension Cords", []string{"Burn"}, nil),
	plib.NewBasicItem("Cell Phone Charger", nil, nil),
	plib.NewBasicItem("Batteries for headlamp", nil, nil),
	plib.NewBasicItem("Sound meter", []string{"Burn"}, nil),
	plib.NewBasicItem("2x modular rigs", []string{"Burn"}, nil),
	plib.NewBasicItem("2x modular power supplies", []string{"Burn"}, nil),
	plib.NewBasicItem("patch cables", []string{"Burn"}, nil),
	plib.NewBasicItem("volca sample", []string{"Burn"}, nil),
	plib.NewBasicItem("retrokits midi cable", []string{"Burn"}, nil),
	plib.NewBasicItem("regular midi cable", []string{"Burn"}, nil),
	plib.NewBasicItem("midi 4x", []string{"Burn"}, nil),
	plib.NewBasicItem("volca power supply", []string{"Burn"}, nil),
	plib.NewBasicItem("octopus power strip", []string{"Burn"}, nil),
	plib.NewBasicItem("DJ headphones", []string{"Burn"}, nil),
	plib.NewBasicItem("long RCA cable", []string{"Burn"}, nil),
	plib.NewBasicItem("1/8 to RCA", []string{"Burn"}, nil),
	plib.NewBasicItem("6x 1/4in adapters", []string{"Burn"}, nil),
	plib.NewBasicItem("DECISION: mixer?", []string{"Burn"}, nil),
	plib.NewBasicItem("beatstep pro", []string{"Burn"}, nil),
	plib.NewBasicItem("beatstep pro power", []string{"Burn"}, nil),
	plib.NewBasicItem("2x beatstep midi adapter", []string{"Burn"}, nil),
}

func init() {
	plib.RegisterItems("Electrical", electrical)
}
