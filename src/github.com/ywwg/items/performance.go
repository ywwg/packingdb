package items

import (
	plib "github.com/ywwg/packinglib"
)

var Performing = []plib.Item{
	plib.NewBasicItem("1x modular rigs", []string{"Performing"}, nil),
	plib.NewBasicItem("1x modular power supplies", []string{"Performing"}, nil),
	plib.NewBasicItem("patch cables", []string{"Performing"}, nil),
	plib.NewBasicItem("modular 3ring", []string{"Performing"}, nil),
	plib.NewBasicItem("2x regular midi cable", []string{"Performing"}, nil),
	plib.NewBasicItem("short midi cable", []string{"Performing"}, nil),
	plib.NewBasicItem("midi 4x", []string{"Performing"}, nil),
	plib.NewBasicItem("octopus power strip", []string{"Performing"}, nil),
	plib.NewBasicItem("DJ headphones", []string{"Performing"}, nil),
	plib.NewBasicItem("long RCA cable", []string{"Performing"}, nil),
	plib.NewBasicItem("RCA stereo female/female", []string{"Performing"}, nil),
	plib.NewBasicItem("1x 1/8 to RCA", []string{"Performing"}, nil),
	plib.NewBasicItem("6x 1/4in adapters", []string{"Performing"}, nil),
	plib.NewBasicItem("stereo 1/4 to 1/8 headphones", []string{"Performing"}, nil),
	plib.NewBasicItem("beatstep pro", []string{"Performing"}, nil),
	plib.NewBasicItem("beatstep pro power", []string{"Performing"}, nil),
	plib.NewBasicItem("2x beatstep midi adapter", []string{"Performing"}, nil),
	plib.NewBasicItem("Zoom recorder", []string{"Performing"}, nil),
	plib.NewBasicItem("RCA to 1/4 for zoom", []string{"Performing"}, nil),
	// TODO: Create a cumulative item for 1/4in adapters
	plib.NewBasicItem("2x 1/4in adapter for zoom", []string{"Performing"}, nil),
	plib.NewBasicItem("microphone", []string{"Performing"}, nil),
}

func init() {
	plib.RegisterItems("Performing", Performing)
}
