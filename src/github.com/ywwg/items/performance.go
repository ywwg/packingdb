package items

import (
	plib "github.com/ywwg/packinglib"
)

var Performing = []plib.Item{
	plib.NewBasicItem("2x modular rigs", []string{"Performing"}, nil),
	plib.NewBasicItem("piece of wood", []string{"Performing"}, nil),
	plib.NewBasicItem("2x modular power supplies", []string{"Performing"}, nil),
	plib.NewBasicItem("patch cables", []string{"Performing"}, nil),
	plib.NewBasicItem("volca sample", []string{"Performing"}, nil),
	plib.NewBasicItem("retrokits midi cable", []string{"Performing"}, nil),
	plib.NewBasicItem("regular midi cable", []string{"Performing"}, nil),
	plib.NewBasicItem("short midi cable", []string{"Performing"}, nil),
	plib.NewBasicItem("midi 4x", []string{"Performing"}, nil),
	plib.NewBasicItem("volca power supply", []string{"Performing"}, nil),
	plib.NewBasicItem("octopus power strip", []string{"Performing"}, nil),
	plib.NewBasicItem("DJ headphones", []string{"Performing"}, nil),
	plib.NewBasicItem("long RCA cable", []string{"Performing"}, nil),
	plib.NewBasicItem("2x 1/8 to RCA", []string{"Performing"}, nil),
	plib.NewBasicItem("6x 1/4in adapters", []string{"Performing"}, nil),
	plib.NewBasicItem("stereo 1/4 to 1/8 for moochers", []string{"Performing"}, nil),
	plib.NewBasicItem("XLR to 1/4", []string{"Performing"}, nil),
	plib.NewBasicItem("beatstep pro", []string{"Performing"}, nil),
	plib.NewBasicItem("beatstep pro power", []string{"Performing"}, nil),
	plib.NewBasicItem("2x beatstep midi adapter", []string{"Performing"}, nil),
	plib.NewBasicItem("laptop keyboard cover", []string{"Performing"}, nil),
	plib.NewBasicItem("microphone", []string{"Performing"}, nil),
	plib.NewBasicItem("XLR cable", []string{"Performing"}, nil),
	plib.NewBasicItem("Sound meter", []string{"Burn"}, nil),
	plib.NewBasicItem("music ear plugs", []string{"Performing", "Partying"}, nil),
}

func init() {
	plib.RegisterItems("Performing", Performing)
}
