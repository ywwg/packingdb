package items

import (
	plib "github.com/ywwg/packinglib"
)

var burnDinner = []plib.Item{
	plib.NewBasicItem("8.5 lbs pulled pork", []string{"Burn"}, nil),
	plib.NewBasicItem("3.5 gallons vegan chili", []string{"Burn"}, nil),
	plib.NewBasicItem("50oz pulled jackfruit", []string{"Burn"}, nil),
	plib.NewBasicItem("2 lbs shredded cheese", []string{"Burn"}, nil),
	plib.NewBasicItem("24oz sour cream", []string{"Burn"}, nil),
	plib.NewBasicItem("sarina's cooler", []string{"Camping"}, nil),
}

var saturdayDinner = []plib.Item{
  plib.NewBasicItem("TODO veggies??? (ask may-lee)", []string{"Burn"}, nil),
	plib.NewBasicItem("owen's cooler", []string{"Burn"}, nil),
}


func init() {
	plib.RegisterItems("Food: Burn Night Dinner", burnDinner)
	plib.RegisterItems("Food: Saturday Dinner Fresh Ingredients", burnDinner)
}
