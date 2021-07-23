package items

import (
	plib "github.com/ywwg/packinglib"
)

var waterStuff = []plib.Item{
	plib.NewConsumableMaxItem("swim suit", 0.25, 2.0, plib.NoUnits, []string{"Swimming", "Lodging"}, nil),
	plib.NewBasicItem("swim towel", []string{"Swimming"}, nil),
	plib.NewBasicItem("drybag", []string{"Swimming", "Dirt"}, nil),
}

var tubing = []plib.Item{
	plib.NewBasicItem("innertube", []string{"Tubing"}, nil),
	plib.NewBasicItem("air pump", []string{"Tubing"}, nil),
}

func init() {
	plib.RegisterItems("Swim", waterStuff)
	plib.RegisterItems("Tubing", tubing)
}
