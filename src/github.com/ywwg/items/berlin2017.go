package items

import (
	plib "github.com/ywwg/packinglib"
)

var berlin2017 = []plib.Item{
	plib.NewBasicItem("Reichstag Ticket", []string{"Berlin2017"}, nil),
	plib.NewBasicItem("Tripadvisor Tickets", []string{"Berlin2017"}, nil),
	plib.NewBasicItem("Lodging information", []string{"Berlin2017"}, nil),
	plib.NewBasicItem("Techno jacket", []string{"Berlin2017"}, nil),
	plib.NewBasicItem("rave pants", []string{"Berlin2017"}, nil),
}

func init() {
	plib.RegisterItems("Berlin2017", berlin2017)
}
