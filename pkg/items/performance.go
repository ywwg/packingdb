package items

import (
	plib "github.com/ywwg/packinglib"
)

var performing = []plib.Item{
	plib.NewBasicItem("ext cord", []string{"Performing"}, nil),
	plib.NewBasicItem("HD25 headphones", []string{"Performing"}, nil),
	plib.NewBasicItem("long RCA cable", []string{"Performing"}, nil),
	plib.NewBasicItem("RCA stereo female/female", []string{"Performing"}, nil),
	plib.NewBasicItem("1x 1/8 to RCA", []string{"PA System", "Projector"}, nil),
	plib.NewBasicItem("stereo 1/4 to 1/8 headphones", []string{"Performing"}, nil),
	plib.NewBasicItem("DJ flight case", []string{"DJing"}, nil),
	plib.NewBasicItem("DJ Controller Power", []string{"DJing"}, nil),
	plib.NewBasicItem("usb-C to usb-B cable", []string{"DJing"}, nil),
	plib.NewBasicItem("RCA to 1/4 for DJ Audio", []string{"Performing"}, nil),
	plib.NewBasicItem("laptop", []string{"DJing", "Karaoke", "Projector"}, nil),
	plib.NewBasicItem("little usb audio", []string{"DJing"}, nil),
	plib.NewBasicItem("MiniMixxx controller", []string{"MiniMixxx"}, nil),
	plib.NewBasicItem("small philips", []string{"MiniMixxx"}, nil),
	plib.NewBasicItem("needlenose pliers", []string{"MiniMixxx"}, nil),
	plib.NewBasicItem("mini bt keyboard", []string{"MiniMixxx"}, nil),
	plib.NewBasicItem("GPD Micro", []string{"MiniMixxx"}, nil),
	plib.NewBasicItem("Huxley stickers", []string{"fancon"}, nil),
	plib.NewBasicItem("2x USB-A to square", []string{"MiniMixxx"}, nil),
	plib.NewBasicItem("Traktor Audio10", []string{"MiniMixxx"}, nil),
	plib.NewBasicItem("Spare USB-C charger", []string{"MiniMixxx"}, nil),
}

func init() {
	plib.RegisterItems("Performing", performing)
}
