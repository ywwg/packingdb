package items

import (
	plib "github.com/ywwg/packinglib"
)

var nevermore = []plib.Item{
	plib.NewBasicItem("Nevermore", []string{"Burn"}, nil),
	plib.NewBasicItem("EZ Up", []string{"Burn"}, nil),
	plib.NewBasicItem("floor tarp", []string{"Burn"}, nil),
	plib.NewBasicItem("3way ext cord", []string{"Burn"}, nil),
	plib.NewBasicItem("UPS", []string{"Burn"}, nil),
	plib.NewBasicItem("airbed (fix it??)", []string{"Burn"}, nil),
	plib.NewBasicItem("2 power strips", []string{"Burn"}, nil),
	plib.NewBasicItem("6 USB power bricks", []string{"Burn"}, nil),
	plib.NewBasicItem("12 CHECKED pis", []string{"Burn"}, nil),
	plib.NewBasicItem("14 usb cords for pis", []string{"Burn"}, nil),
	plib.NewBasicItem("CONFIGURED wifi point", []string{"Burn"}, nil),
	plib.NewBasicItem("power brick for wifi point", []string{"Burn"}, nil),
	plib.NewBasicItem("12 speakers with cords", []string{"Burn"}, nil),
	plib.NewBasicItem("2 copper light strings", []string{"Burn"}, nil),
	plib.NewBasicItem("6 wooden bracket things", []string{"Burn"}, nil),
	plib.NewBasicItem("round table", []string{"Burn"}, nil),
	plib.NewBasicItem("olpc", []string{"Burn"}, nil),
	plib.NewBasicItem("midi fighter w/cable", []string{"Burn"}, nil),
	plib.NewBasicItem("nevermore binder", []string{"Burn"}, nil),
	plib.NewBasicItem("solar string lights", []string{"Burn"}, nil),
	plib.NewBasicItem("nevermore sign", []string{"Burn"}, nil),
	plib.NewBasicItem("rope to hold sign", []string{"Burn"}, nil),
}

func init() {
	plib.RegisterItems("Nevermore", nevermore)
}
