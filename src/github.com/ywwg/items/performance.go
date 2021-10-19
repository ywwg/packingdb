package items

import (
	plib "github.com/ywwg/packinglib"
)

var performing = []plib.Item{
	plib.NewBasicItem("1x modular rigs", []string{"Modular"}, nil),
	plib.NewBasicItem("1x modular power supplies", []string{"Modular"}, nil),
	plib.NewBasicItem("patch cables", []string{"Modular"}, nil),
	plib.NewBasicItem("modular 3ring", []string{"Modular"}, nil),
	plib.NewBasicItem("2x regular midi cable", []string{"Modular"}, nil),
	plib.NewBasicItem("short midi cable", []string{"Modular"}, nil),
	plib.NewBasicItem("midi 4x", []string{"Modular"}, nil),
	plib.NewBasicItem("beatstep pro", []string{"Modular"}, nil),
	plib.NewBasicItem("beatstep pro power", []string{"Modular"}, nil),
	plib.NewBasicItem("2x beatstep midi adapter", []string{"Modular"}, nil),
	plib.NewBasicItem("microsd card reader", []string{"Modular"}, nil),
	plib.NewBasicItem("octopus power strip", []string{"Performing"}, nil),
	plib.NewBasicItem("DJ headphones", []string{"Performing"}, nil),
	plib.NewBasicItem("long RCA cable", []string{"Performing"}, nil),
	plib.NewBasicItem("RCA stereo female/female", []string{"Performing"}, nil),
	plib.NewBasicItem("1x 1/8 to RCA", []string{"Modular", "PA System", "Projector"}, nil),
	plib.NewBasicItem("6x 1/4in adapters", []string{"Modular"}, nil),
	plib.NewBasicItem("stereo 1/4 to 1/8 headphones", []string{"Performing"}, nil),
	plib.NewBasicItem("Zoom recorder", []string{"Performing"}, nil),
	plib.NewBasicItem("Zoom power", []string{"Performing"}, nil),
	plib.NewBasicItem("RCA to 1/4 for zoom", []string{"Performing"}, nil),
	// TODO: Create a cumulative item for 1/4in adapters
	plib.NewBasicItem("2x 1/4in adapter for zoom", []string{"Performing"}, nil),
	plib.NewBasicItem("microphone", []string{"Performing"}, nil),
	plib.NewBasicItem("DJ Controller", []string{"DJing"}, nil),
	plib.NewBasicItem("DJ Controller Power", []string{"DJing"}, nil),
	plib.NewBasicItem("bag for DJ controller", []string{"DJing"}, nil),
	plib.NewBasicItem("usb-C to usb-B cable", []string{"DJing"}, nil),
	plib.NewBasicItem("RCA to 1/4 for DJ audio", []string{"Performing"}, nil),
	plib.NewBasicItem("laptop", []string{"DJing", "Karaoke", "Projector"}, nil),
	plib.NewBasicItem("keyboard cover", []string{"DJing"}, nil),
	plib.NewBasicItem("laptop stand", []string{"DJing"}, nil),
	plib.NewBasicItem("little usb audio", []string{"DJing"}, nil),
	plib.NewBasicItem("MiniMixxx", []string{"MiniMixxx"}, nil),
	plib.NewBasicItem("small philips", []string{"MiniMixxx"}, nil),
	plib.NewBasicItem("needlenose pliers", []string{"MiniMixxx"}, nil),
	plib.NewBasicItem("mini bt keyboard", []string{"MiniMixxx"}, nil),
	plib.NewBasicItem("MiniMixxx", []string{"MiniMixxx"}, nil),
	plib.NewBasicItem("Mixxx stickers", []string{"Performing"}, nil),
	plib.NewBasicItem("2x USB-A to square", []string{"MiniMixxx"}, nil),
	plib.NewBasicItem("Traktor Audio10", []string{"MiniMixxx"}, nil),
	plib.NewBasicItem("Spare USB-C charger", []string{"MiniMixxx"}, nil),
}

func init() {
	plib.RegisterItems("Performing", performing)
}
