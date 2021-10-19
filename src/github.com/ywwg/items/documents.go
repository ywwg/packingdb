package items

import (
	plib "github.com/ywwg/packinglib"
)

var documents = []plib.Item{
	plib.NewBasicItem("Wallet", nil, nil),
	plib.NewBasicItem("Passport", []string{"International"}, nil),
	plib.NewBasicItem("Global Entry Card", []string{"Flight"}, nil),
	plib.NewBasicItem("Event Ticket", []string{"PaidEvent"}, nil),
	plib.NewBasicItem("Print Wiki", []string{"BigTrip"}, nil),
	plib.NewBasicItem("Waiver", []string{"Burn"}, nil),
	plib.NewBasicItem("Transportation Tickets", []string{"PaidTravel", "Flight"}, nil),
	plib.NewBasicItem("MTA Metrocard", []string{"NYC"}, nil),
	plib.NewBasicItem("AAA card", []string{"Tiny House", "CarRide"}, nil),
	plib.NewBasicItem("VAX card", []string{"Vax Proof"}, nil),
	plib.NewBasicItem("COVID result proof", []string{"Vax Proof"}, nil),
}

func init() {
	plib.RegisterItems("Documents", documents)
}
