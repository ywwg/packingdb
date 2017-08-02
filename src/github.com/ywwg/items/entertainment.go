package items

import (
	plib "github.com/ywwg/packinglib"
)

var entertainment = []plib.Item{
	plib.NewConsumableMaxItem("books", 0.5, 3.0, plib.NoUnits, nil, nil),
	plib.NewBasicItem("earbuds", nil, nil),
	plib.NewBasicItem("laptop and charger", nil, nil),
	plib.NewBasicItem("USB-C cable", nil, []string{"Tiny House"}),
	plib.NewBasicItem("jambox and charger", nil, []string{"International"}),
	plib.NewCustomConsumableItem("tv / movie", func(nights int, props plib.PropertySet) float64 {
		if _, ok := props["Tiny House"]; ok {
			return float64(nights)
		}
		return 3.0
	}, plib.NoUnits, []string{"Tiny House", "International"}, nil),
	plib.NewBasicItem("music ear plugs", []string{"Performing", "Partying"}, nil),
}

func init() {
	plib.RegisterItems("Entertainment", entertainment)
}
