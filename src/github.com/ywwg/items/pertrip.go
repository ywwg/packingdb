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

// var virtuality2017 = []plib.Item{
// 	plib.NewBasicItem("camp platewear kit", []string{"Retreat2017"}, nil),
// 	plib.NewBasicItem("canned beer", []string{"Retreat2017"}, nil),
// }

var firefly2018 = []plib.Item{
	plib.NewBasicItem("Spreadsheet items", []string{"Firefly2018"}, nil),
	plib.NewBasicItem("2 5gal water", []string{"Firefly2018"}, nil),
	plib.NewBasicItem("2 doz eggs", []string{"Firefly2018"}, nil),
	plib.NewBasicItem("3 pkgs jackfruit", []string{"Firefly2018"}, nil),
	plib.NewBasicItem("2 cans refried beans", []string{"Firefly2018"}, nil),
	plib.NewBasicItem("tupperware cut peppers", []string{"Firefly2018"}, nil),
	plib.NewBasicItem("1 pkg cheese", []string{"Firefly2018"}, nil),
	plib.NewBasicItem("1 jar jalapenos", []string{"Firefly2018"}, nil),
	plib.NewBasicItem("2 bottle hot sauce", []string{"Firefly2018"}, nil),
	plib.NewBasicItem("2 jar salsa", []string{"Firefly2018"}, nil),
	plib.NewBasicItem("corn tortillas", []string{"Firefly2018"}, nil),
	plib.NewBasicItem("flour tortillas", []string{"Firefly2018"}, nil),
}

var pride = []plib.Item{
	plib.NewBasicItem("pride shirt", []string{"PortlandPride"}, nil),
}

func init() {
	plib.RegisterItems("Berlin2017", berlin2017)
	// plib.RegisterItems("Retreat2017", virtuality2017)
	plib.RegisterItems("Firefly2018", firefly2018)
	plib.RegisterProperty(plib.Property("PortlandPride"))
	plib.RegisterItems("PortlandPride", pride)
}
