package items

import (
	plib "github.com/ywwg/packingdb/pkg/packinglib"
)

var documents = []*plib.Item{
	plib.NewItem("Wallet", nil, nil),
	plib.NewItem("Keys (incl yubi)", nil, nil),
	plib.NewItem("Passport", []string{"International"}, nil),
	plib.NewItem("Event Tickets", []string{"PaidEvent"}, nil),
	plib.NewItem("Print Reservations", []string{"BigTrip"}, nil),
	plib.NewItem("Transportation Tickets", []string{"PaidTravel", "Flight"}, nil),
	plib.NewItem("AAA card", []string{"Tiny House", "CarRide"}, nil),
	plib.NewItem("Tile in suitcase", []string{"PaidTravel"}, []string{"NoCheckedLuggage"}),
	plib.NewItem("COVID tests", []string{"Vax Proof"}, nil).Consumable(0.5),
}
