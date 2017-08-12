package items

import (
	plib "github.com/ywwg/packinglib"
)

var entertainment = []plib.Item{
	plib.NewConsumableMaxItem("books", 0.5, 3.0, plib.NoUnits, nil, nil),
	plib.NewBasicItem("guidebooks", []string{"Big Trip"}, nil),
	plib.NewBasicItem("earbuds", nil, nil),
	plib.NewBasicItem("laptop and charger", nil, nil),
	plib.NewBasicItem("USB-C cable", nil, []string{"Tiny House"}),
	plib.NewBasicItem("jambox and charger", nil, []string{"International"}),
	plib.NewBasicItem("1/8 stereo cable", nil, []string{"International"}),
	plib.NewCustomConsumableItem("tv / movie", func(nights int, props plib.PropertySet) float64 {
		if _, ok := props["Tiny House"]; ok {
			return float64(nights)
		}
		if _, ok := props["Flight"]; ok {
			// Should be enough for two plane flights and random nights.
			return 4.0
		}
		return 0.0
	}, plib.NoUnits, []string{"Tiny House", "International"}, nil),
	plib.NewBasicItem("music ear plugs", []string{"Performing", "Partying"}, nil),
}

func init() {
	plib.RegisterItems("Entertainment", entertainment)
}
