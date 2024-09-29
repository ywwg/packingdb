package items

import (
	plib "github.com/ywwg/packingdb/pkg/packinglib"
)

var toiletries = []*plib.Item{
	plib.NewItem("deoderant", nil, nil),
	plib.NewItem("wet wipes", []string{"Camping", "Tiny House", "Playtime"}, nil),
	plib.NewItem("hand sanitizer", []string{"Camping", "Dirt"}, nil),
	plib.NewItem("2 kitchen towels", []string{"Tiny House"}, nil),
	plib.NewItem("toothpaste", nil, nil),
	plib.NewItem("toothbrush", nil, nil),
	plib.NewItem("shampoo", nil, []string{"Burn", "NoCheckedLuggage", "HasToiletries"}),
	plib.NewItem("soap", nil, []string{"Burn", "HasToiletries"}),
	plib.NewItem("flossers", nil, []string{"HasToiletries"}),
	plib.NewItem("nail clippers", nil, nil),
	plib.NewItem("cute clippers", nil, nil),
	plib.NewItem("tweezers", nil, nil),
	plib.NewItem("tissues", nil, nil),
	plib.NewItem("shaving Cream", nil, []string{"Burn"}),
	plib.NewItem("aftershave", nil, nil),
	plib.NewItem("hair stuff", nil, []string{"Burn", "NoCheckedLuggage"}),
	plib.NewItem("hair dye", nil, []string{"Burn"}).Custom(func(_ float64, nights int, props plib.PropertySet) float64 {
		if nights >= 3 {
			return 1
		}
		return 0
	}),
	plib.NewItem("travel hair stuff", []string{"NoCheckedLuggage"}, nil),
	plib.NewItem("electric razor", []string{"Burn", "Tiny House"}, nil),
	plib.NewItem("razor", nil, []string{"Burn"}),
	plib.NewItem("spray sunscreen", []string{"Camping", "Bright", "Hiking", "Boat"}, []string{"Flight"}).TemperatureRange(70, 120),
	plib.NewItem("regular sunscreen", []string{"Camping", "Bright", "Hiking", "Boat"}, nil).TemperatureRange(70, 120),
	plib.NewItem("lip sunscreen", []string{"Camping", "Bright", "Hiking", "Boat"}, nil).TemperatureRange(70, 120),
	plib.NewItem("regular glasses", nil, nil),
	plib.NewItem("glasses cleaner", nil, nil),
	plib.NewItem("sunglasses", nil, nil),
	plib.NewItem("oakleys", nil, []string{"Burn"}),
	plib.NewItem("contact Fluid", nil, []string{"Burn"}),
	plib.NewItem("contacts", nil, []string{"Burn"}).Units("pairs").Consumable(1.0),
	plib.NewItem("epi pens", nil, nil),
	plib.NewItem("Dramamine", []string{"CarRide", "Boat"}, nil),
	plib.NewItem("anti-nausea patches", []string{"Boat"}, nil),
}
