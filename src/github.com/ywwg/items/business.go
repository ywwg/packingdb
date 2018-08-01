package items

import (
	plib "github.com/ywwg/packinglib"
)

var business = []plib.Item{
	plib.NewConsumableItem("business shirts", 0.75, plib.NoUnits, []string{"Business"}, nil),
	plib.NewBasicItem("work laptop", []string{"Business"}, nil),
	plib.NewBasicItem("work laptop power", []string{"Business"}, nil),
	plib.NewBasicItem("badge", []string{"Business"}, nil),
	plib.NewBasicItem("gcard", []string{"Business"}, nil),
}

func init() {
	plib.RegisterItems("Business", business)
}
