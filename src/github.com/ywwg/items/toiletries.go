package items

import (
	plib "github.com/ywwg/packinglib"
)

var toiletries = []plib.Item{
	plib.NewBasicItem("Deoderant", nil, nil),
	plib.NewBasicItem("Wet Wipes", []string{"Burn"}, nil),
	plib.NewBasicItem("Toothpaste", nil, []string{"HasToiletries"}),
	plib.NewBasicItem("Toothbrush", nil, nil),
	plib.NewBasicItem("Shampoo", nil, []string{"Burn", "NoCheckedLuggage", "HasToiletries"}),
	plib.NewBasicItem("Soap", nil, []string{"Burn", "HasToiletries"}),
	plib.NewBasicItem("Flossers", nil, []string{"HasToiletries"}),
	plib.NewBasicItem("Nail Clippers", nil, nil),
	plib.NewBasicItem("Cute Clippers", nil, nil),
	plib.NewBasicItem("Shaving Cream", nil, []string{"Burn"}),
	plib.NewBasicItem("Aftershave", nil, []string{"Burn"}),
	plib.NewBasicItem("Razor", nil, []string{"Burn"}),
	plib.NewBasicItem("Spray Sunscreen", []string{"Camping", "Bright"}, []string{"NoCheckedLuggage"}),
	plib.NewBasicItem("Regular Sunscreen", []string{"Camping", "Bright"}, nil),
	plib.NewBasicItem("Lip Sunscreen", []string{"Camping", "Bright"}, nil),
	plib.NewBasicItem("Regular glasses", nil, nil),
	plib.NewBasicItem("Sunglasses", nil, nil),
	plib.NewBasicItem("Oakleys", []string{"Contacts"}, nil),
	plib.NewBasicItem("Contact Fluid", []string{"Contacts"}, nil),
	plib.NewConsumableItem("Contacts", 2.0, plib.NoUnits, []string{"Contacts"}, nil),
}

func init() {
	plib.RegisterItems("Toiletries", toiletries)
}
