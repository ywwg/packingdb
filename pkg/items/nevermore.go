package items

import (
	plib "github.com/ywwg/packingdb/pkg/packinglib"
)

var nevermore = []plib.Item{
	plib.NewBasicItem("Nevermore", []string{"Nevermore"}, nil),
	plib.NewBasicItem("EZ Up", []string{"Nevermore"}, nil),
	plib.NewBasicItem("floor tarp", []string{"Nevermore"}, nil),
	plib.NewBasicItem("3way ext cord", []string{"Nevermore"}, nil),
	plib.NewBasicItem("UPS", []string{"Nevermore"}, nil),
	plib.NewBasicItem("airbed", []string{"Nevermore"}, nil),
	plib.NewBasicItem("2 power strips", []string{"Nevermore"}, nil),
	plib.NewBasicItem("4 USB power bricks", []string{"Nevermore"}, nil),
	plib.NewBasicItem("12 CHECKED pis", []string{"Nevermore"}, nil),
	plib.NewBasicItem("14 usb cords for pis (2 extra)", []string{"Nevermore"}, nil),
	plib.NewBasicItem("CONFIGURED wifi point", []string{"Nevermore"}, nil),
	plib.NewBasicItem("power brick for wifi point", []string{"Nevermore"}, nil),
	plib.NewBasicItem("2 copper light strings", []string{"Nevermore"}, nil),
	plib.NewBasicItem("6 wooden bracket things", []string{"Nevermore"}, nil),
	plib.NewBasicItem("round table", []string{"Nevermore"}, nil),
	plib.NewBasicItem("4 screws for table", []string{"Nevermore"}, nil),
	plib.NewBasicItem("drill to assemble table", []string{"Nevermore"}, nil),
	plib.NewBasicItem("olpc", []string{"Nevermore"}, nil),
	plib.NewBasicItem("midi fighter w/cable", []string{"Nevermore"}, nil),
	plib.NewBasicItem("nevermore binder", []string{"Nevermore"}, nil),
	plib.NewBasicItem("tea lights", []string{"Nevermore"}, nil),
	plib.NewBasicItem("CR2032 batts for tea lights", []string{"Nevermore"}, nil),
	plib.NewBasicItem("nevermore sign", []string{"Nevermore"}, nil),
	plib.NewBasicItem("rope to hold sign", []string{"Nevermore"}, nil),
}

func init() {
	plib.RegisterItems("Nevermore", nevermore)
}
