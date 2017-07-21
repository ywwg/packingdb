package items

import (
	plib "github.com/ywwg/packinglib"
)

var waterStuff = []plib.Item{
	plib.NewBasicItem("froggies", []string{"Swimming"}, nil),
	plib.NewConsumableTemperatureItem("swim suit", 0.25, plib.NoUnits, 70, 120, nil, nil),
	plib.NewBasicItem("swim towel", []string{"Swimming"}, nil),
}

func init() {
	plib.RegisterItems("Swim", waterStuff)
}
