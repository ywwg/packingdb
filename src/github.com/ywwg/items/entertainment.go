package items

import (
	plib "github.com/ywwg/packinglib"
)

var entertainment = []plib.Item{
	plib.NewConsumableMaxItem("books", 1.0, 3.0, plib.NoUnits, nil, nil),
	plib.NewBasicItem("35mm camera and/or polaroid", nil, nil),
	plib.NewBasicItem("guidebooks", []string{"BigTrip"}, nil),
	plib.NewBasicItem("earbuds", nil, nil),
	plib.NewBasicItem("laptop and charger", nil, []string{"Insecure"}),
	plib.NewBasicItem("USB-C breakout", nil, []string{"Insecure"}),
	plib.NewBasicItem("HDMI cable", []string{"Lodging", "Projector"}, []string{"Insecure"}),
	plib.NewBasicItem("USB-A-C cable", nil, []string{"Tiny House"}),
	plib.NewBasicItem("USB-C-C cable", nil, []string{"Tiny House"}),
	plib.NewBasicItem("jambox and charger", []string{"Speaker", "Projector"}, []string{"International"}),
	plib.NewBasicItem("1/8 stereo cable", []string{"Speaker"}, []string{"International", "Tiny House"}),
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
	plib.NewBasicItem("music ear plugs", []string{"Performing", "Partying", "Loud"}, nil),
	plib.NewBasicItem("Mixxx stickers", []string{"Performing"}, nil),
	plib.NewBasicItem("Projector", []string{"Projector"}, nil),
	plib.NewBasicItem("Screen", []string{"Projector"}, nil),
	plib.NewBasicItem("JBL Speakers", []string{"Projector", "PA System"}, nil),
	plib.NewBasicItem("Speaker Stands", []string{"Projector", "PA System"}, nil),
	plib.NewBasicItem("Mixer", []string{"Projector", "PA System"}, nil),
	plib.NewBasicItem("long XLR cables", []string{"Projector", "PA System"}, nil),
	plib.NewBasicItem("Karaoke hard drive", []string{"Karaoke"}, nil),
	plib.NewBasicItem("Extra Mic and XLR", []string{"Karaoke"}, nil),
	plib.NewBasicItem("Party Laser", []string{"Karaoke"}, nil),
}

func init() {
	plib.RegisterItems("Entertainment", entertainment)
}
