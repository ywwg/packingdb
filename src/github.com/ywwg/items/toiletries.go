package items

import (
	plib "github.com/ywwg/packinglib"
)

var toiletries = []plib.Item{
	plib.NewBasicItem("deoderant", nil, nil),
	plib.NewBasicItem("wet wipes", []string{"Camping", "Tiny House"}, nil),
	plib.NewBasicItem("towel", []string{"Tiny House"}, nil),
	plib.NewBasicItem("toothpaste", nil, nil),
	plib.NewBasicItem("toothbrush", nil, nil),
	plib.NewBasicItem("shampoo", nil, []string{"Burn", "NoCheckedLuggage", "HasToiletries"}),
	plib.NewBasicItem("soap", nil, []string{"Burn", "HasToiletries"}),
	plib.NewBasicItem("flossers", nil, []string{"HasToiletries"}),
	plib.NewBasicItem("nail clippers", nil, nil),
	plib.NewBasicItem("cute clippers", nil, nil),
	plib.NewBasicItem("tissues", nil, nil),
	plib.NewBasicItem("shaving Cream", nil, []string{"Burn"}),
	plib.NewBasicItem("aftershave", nil, []string{"Burn"}),
	plib.NewBasicItem("hair stuff", nil, []string{"Burn"}),
	plib.NewBasicItem("electric razor", []string{"Burn"}, nil),
	plib.NewBasicItem("razor", nil, []string{"Burn"}),
	plib.NewBasicItem("spray sunscreen", []string{"Camping", "Bright"}, []string{"Flight"}),
	plib.NewBasicItem("regular sunscreen", []string{"Camping", "Bright"}, nil),
	plib.NewBasicItem("lip sunscreen", []string{"Camping", "Bright"}, nil),
	plib.NewBasicItem("regular glasses", nil, nil),
	plib.NewBasicItem("glasses cleaner", nil, nil),
	plib.NewBasicItem("sunglasses", nil, nil),
	plib.NewBasicItem("oakleys", nil, []string{"Burn"}),
	plib.NewBasicItem("contact Fluid", nil, []string{"Burn"}),
	plib.NewConsumableItem("contacts", 1.0, "pairs", nil, []string{"Burn"}),
}

func init() {
	plib.RegisterItems("Toiletries", toiletries)
}
