package packinglib

var toiletries = []Item{
	NewBasicItem("Deoderant", nil),
	NewBasicItem("Wet Wipes", PropertySet{Burn: true}),
	NewBasicItem("Toothpaste", nil),
	NewBasicItem("Toothbrush", nil),
	NewBasicItem("Shampoo", PropertySet{Burn: false, NoCheckedLuggage: false}),
	NewBasicItem("Soap", PropertySet{Burn: false}),
	NewBasicItem("Flossers", nil),
	NewBasicItem("Nail Clippers", nil),
	NewBasicItem("Cute Clippers", nil),
	NewBasicItem("Flossers", nil),
	NewBasicItem("Shaving Cream", PropertySet{Burn: false}),
	NewBasicItem("Razor", PropertySet{Burn: false}),
	NewBasicItem("Spray Sunscreen", PropertySet{Camping: true, NoCheckedLuggage: false}),
	NewBasicItem("Regular Sunscreen", PropertySet{Camping: true}),
	NewBasicItem("Lip Sunscreen", PropertySet{Camping: true}),
}

func init() {
	RegisterItems("Toiletries", toiletries)
}
