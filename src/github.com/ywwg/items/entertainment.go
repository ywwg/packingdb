package items

import (
	plib "github.com/ywwg/packinglib"
)

var entertainment = []plib.Item{
	plib.NewConsumableItem("books", 0.5, plib.NoUnits, nil, nil),
	plib.NewBasicItem("earbuds", nil, nil),
	plib.NewBasicItem("laptop and charger", nil, nil),
	plib.NewBasicItem("USB-C cable", nil, nil),
}

func init() {
	plib.RegisterItems("Entertainment", entertainment)
}
