package items

import (
	plib "github.com/ywwg/packingdb/pkg/packinglib"
)

var business = []*plib.Item{
	plib.NewItem("business shirts", []string{"Business"}, nil).
		Consumable(0.50).
		Max(3.0),
	plib.NewItem("work laptop", []string{"Business"}, nil),
	plib.NewItem("work laptop power", []string{"Business"}, nil),
	plib.NewItem("headset earbuds", []string{"Business"}, nil),
	plib.NewItem("UGREEN power/headset USB-C", []string{"Business"}, nil),
}
