package items

import (
	plib "github.com/ywwg/packinglib"
)

var documents = []plib.Item{
	plib.NewBasicItem("Wallet", nil, nil),
	plib.NewBasicItem("Passport", []string{"International"}, nil),
	plib.NewBasicItem("Event Tickets", []string{"PaidEvent"}, nil),
	plib.NewBasicItem("Print Reservations", []string{"BigTrip"}, nil),
	plib.NewBasicItem("Transportation Tickets", []string{"PaidTravel", "Flight"}, nil),
	plib.NewBasicItem("AAA card", []string{"Tiny House", "CarRide"}, nil),
	plib.NewBasicItem("Tile in suitcase", []string{"PaidTravel"}, []string{"NoCheckedLuggage"}),
	plib.NewConsumableItem("COVID tests", 0.5, plib.NoUnits, []string{"Vax Proof"}, nil),
}

func init() {
	plib.RegisterItems("Documents", documents)
}
