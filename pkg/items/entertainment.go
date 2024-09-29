package items

import (
	plib "github.com/ywwg/packingdb/pkg/packinglib"
)

var entertainment = []*plib.Item{
	plib.NewItem("books", nil, nil).Consumable(1.0).Max(3.0),
	plib.NewItem("35mm camera and/or polaroid", nil, nil),
	plib.NewItem("guidebooks", []string{"BigTrip"}, nil),
	plib.NewItem("earbuds", nil, nil),
	plib.NewItem("laptop and charger", nil, []string{"Insecure"}),
	plib.NewItem("USB-C breakout", nil, []string{"Insecure"}),
	plib.NewItem("HDMI cable", []string{"Lodging", "Projector"}, []string{"Insecure"}),
	plib.NewItem("USB-C-C cable", nil, []string{"Tiny House"}),
	plib.NewItem("USB-A-C cable", []string{"International", "Flight", "Lodging"}, nil),
	plib.NewItem("bluetooth speaker and charger", []string{"Speaker", "Projector"}, []string{"International"}),
	plib.NewItem("1/8 stereo cable", []string{"Speaker"}, []string{"International", "Tiny House"}),
	plib.NewItem("tv / movie", []string{"Tiny House", "International"}, nil).Custom(func(_ float64, nights int, props plib.PropertySet) float64 {
		if _, ok := props["Tiny House"]; ok {
			return float64(nights)
		}
		if _, ok := props["Flight"]; ok {
			// Should be enough for two plane flights and random nights.
			return 4.0
		}
		return 0.0
	}),
	plib.NewItem("music ear plugs", []string{"Performing", "Partying", "Loud"}, nil),
	plib.NewItem("Projector", []string{"Projector"}, nil),
	plib.NewItem("Screen", []string{"Projector"}, nil),
	plib.NewItem("JBL Speakers", []string{"Projector", "PA System"}, nil),
	plib.NewItem("Speaker Stands", []string{"Projector", "PA System"}, nil),
	plib.NewItem("Mixer", []string{"Projector", "PA System"}, nil),
	plib.NewItem("long XLR cables", []string{"Projector", "PA System"}, nil),
	plib.NewItem("Karaoke hard drive", []string{"Karaoke"}, nil),
	plib.NewItem("Extra Mic and XLR", []string{"Karaoke"}, nil),
	plib.NewItem("Party Laser", []string{"Karaoke"}, nil),
}
