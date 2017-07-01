package items

import (
	plib "github.com/ywwg/packinglib"
)

var toiletries = []plib.Item{
	plib.NewBasicItem("Deoderant", nil, nil),
	plib.NewBasicItem("Wet Wipes", []string{"Burn"}, nil),
	plib.NewBasicItem("Toothpaste", nil, nil),
	plib.NewBasicItem("Toothbrush", nil, nil),
	plib.NewBasicItem("Shampoo", nil, []string{"Burn", "NoCheckedLuggage"}),
	plib.NewBasicItem("Soap", nil, []string{"Burn"}),
	plib.NewBasicItem("wet wipes", []string{"Burn"}, nil),
	plib.NewBasicItem("Flossers", nil, nil),
	plib.NewBasicItem("Nail Clippers", nil, nil),
	plib.NewBasicItem("Cute Clippers", nil, nil),
	plib.NewBasicItem("Shaving Cream", nil, []string{"Burn"}),
	plib.NewBasicItem("Razor", nil, []string{"Burn"}),
	plib.NewBasicItem("Spray Sunscreen", []string{"Camping"}, []string{"NoCheckedLuggage"}),
	plib.NewBasicItem("Regular Sunscreen", []string{"Camping"}, nil),
	plib.NewBasicItem("Lip Sunscreen", []string{"Camping"}, nil),
}

func init() {
	plib.RegisterItems("Toiletries", toiletries)
}
