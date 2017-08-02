package items

import (
	plib "github.com/ywwg/packinglib"
)

var entertainment = []plib.Item{
	plib.NewConsumableItem("books", 0.5, plib.NoUnits, nil, nil),
	plib.NewBasicItem("earbuds", nil, nil),
	plib.NewBasicItem("laptop and charger", nil, nil),
	plib.NewBasicItem("USB-C cable", nil, []string{"Tiny House"}),
	plib.NewBasicItem("jambox and charger", nil, []string{"International"}),
	plib.NewConsumableItem("tv / movie", 1.0, plib.NoUnits, []string{"Tiny House", "International"}, nil),
	plib.NewBasicItem("music ear plugs", []string{"Performing", "Partying"}, nil),
}

func init() {
	plib.RegisterItems("Entertainment", entertainment)
}
