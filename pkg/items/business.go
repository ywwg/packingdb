package items

import (
	plib "github.com/ywwg/packinglib"
)

var business = []plib.Item{
	plib.NewConsumableMaxItem("business shirts", 0.50, 3.0, plib.NoUnits, []string{"Business"}, nil),
	plib.NewBasicItem("work laptop", []string{"Business"}, nil),
	plib.NewBasicItem("work laptop power", []string{"Business"}, nil),
}

func init() {
	plib.RegisterItems("Business", business)
}
