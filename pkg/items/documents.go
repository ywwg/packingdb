package items

import (
	plib "github.com/ywwg/packingdb/pkg/packinglib"
)

var documents = []*plib.Item{
	plib.NewItem("Wallet", nil, nil),
	plib.NewItem("Passport", []string{"International"}, nil),
	plib.NewItem("Event Tickets", []string{"PaidEvent"}, nil),
	plib.NewItem("Print Reservations", []string{"BigTrip"}, nil),
	plib.NewItem("Transportation Tickets", []string{"PaidTravel", "Flight"}, nil),
	plib.NewItem("MTA Metrocard", []string{"NYC"}, nil),
	plib.NewItem("AAA card", []string{"Tiny House", "CarRide"}, nil),
	plib.NewItem("VAX card", []string{"Vax Proof"}, nil),
	plib.NewItem("Tile in suitcase", []string{"PaidTravel"}, []string{"NoCheckedLuggage"}),
	plib.NewConsumableItem("COVID tests", 0.5, plib.NoUnits, []string{"Vax Proof"}, nil),
}

func init() {
	plib.RegisterItems("Documents", documents)
}
