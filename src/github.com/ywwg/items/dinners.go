package items

import (
	plib "github.com/ywwg/packinglib"
)

var burnDinner = []plib.Item{
	plib.NewBasicItem("8.5 lbs pulled pork", []string{"Firefly2017"}, nil),
	plib.NewBasicItem("3.5 gallons vegan chili", []string{"Firefly2017"}, nil),
	plib.NewBasicItem("50oz pulled jackfruit", []string{"Firefly2017"}, nil),
	plib.NewBasicItem("2 lbs shredded cheese", []string{"Firefly2017"}, nil),
	plib.NewBasicItem("24oz sour cream", []string{"Firefly2017"}, nil),
	plib.NewBasicItem("sarina's cooler", []string{"Firefly2017"}, nil),
}

var saturdayDinner = []plib.Item{
	plib.NewBasicItem("TODO veggies??? (ask may-lee)", []string{"Firefly2017"}, nil),
	plib.NewBasicItem("owen's cooler", []string{"Firefly2017"}, nil),
}

func init() {
	plib.RegisterProperty("Firefly2017")
	plib.RegisterItems("Food: Burn Night Dinner", burnDinner)
	plib.RegisterItems("Food: Saturday Dinner Fresh Ingredients", saturdayDinner)
}
