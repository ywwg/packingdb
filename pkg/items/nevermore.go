package items

import (
	plib "github.com/ywwg/packingdb/pkg/packinglib"
)

var nevermore = []*plib.Item{
	plib.NewItem("Nevermore", []string{"Nevermore"}, nil),
	plib.NewItem("EZ Up", []string{"Nevermore", "Burn"}, nil),
	plib.NewItem("floor tarp", []string{"Nevermore"}, nil),
	plib.NewItem("3way ext cord", []string{"Nevermore"}, nil),
	plib.NewItem("UPS", []string{"Nevermore"}, nil),
	plib.NewItem("airbed", []string{"Nevermore"}, nil),
	plib.NewItem("2 power strips", []string{"Nevermore"}, nil),
	plib.NewItem("4 USB power bricks", []string{"Nevermore"}, nil),
	plib.NewItem("12 CHECKED pis", []string{"Nevermore"}, nil),
	plib.NewItem("14 usb cords for pis (2 extra)", []string{"Nevermore"}, nil),
	plib.NewItem("CONFIGURED wifi point", []string{"Nevermore"}, nil),
	plib.NewItem("power brick for wifi point", []string{"Nevermore"}, nil),
	plib.NewItem("2 copper light strings", []string{"Nevermore"}, nil),
	plib.NewItem("6 wooden bracket things", []string{"Nevermore"}, nil),
	plib.NewItem("round table", []string{"Nevermore"}, nil),
	plib.NewItem("4 screws for table", []string{"Nevermore"}, nil),
	plib.NewItem("drill to assemble table", []string{"Nevermore"}, nil),
	plib.NewItem("olpc", []string{"Nevermore"}, nil),
	plib.NewItem("midi fighter w/cable", []string{"Nevermore"}, nil),
	plib.NewItem("nevermore binder", []string{"Nevermore"}, nil),
	plib.NewItem("tea lights", []string{"Nevermore"}, nil),
	plib.NewItem("CR2032 batts for tea lights", []string{"Nevermore"}, nil),
	plib.NewItem("nevermore sign", []string{"Nevermore"}, nil),
	plib.NewItem("rope to hold sign", []string{"Nevermore"}, nil),
}

func init() {
	plib.RegisterItems("Nevermore", nevermore)
}
