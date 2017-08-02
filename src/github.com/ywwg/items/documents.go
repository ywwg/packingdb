package items

import (
	plib "github.com/ywwg/packinglib"
)

var documents = []plib.Item{
	plib.NewBasicItem("Wallet", nil, nil),
	plib.NewBasicItem("Passport", []string{"International"}, nil),
	plib.NewBasicItem("Reichstag Ticket", []string{"Berlin2017"}, nil),
	plib.NewBasicItem("Tripadvisor Tickets", []string{"Berlin2017"}, nil),
	plib.NewBasicItem("Event Ticket", []string{"PaidEvent"}, nil),
}

func init() {
	plib.RegisterItems("Documents", documents)
}
