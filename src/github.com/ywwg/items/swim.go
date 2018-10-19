package items

import (
	plib "github.com/ywwg/packinglib"
)

var waterStuff = []plib.Item{
	plib.NewBasicItem("froggies", []string{"Swimming"}, nil),
	plib.NewConsumableMaxItem("swim suit", 0.25, 2.0, plib.NoUnits, []string{"Swimming", "Lodging"}, nil),
	plib.NewBasicItem("swim towel", []string{"Swimming"}, nil),
	plib.NewBasicItem("plastic bag for wet things", []string{"Swimming", "Dirt"}, nil),
}

func init() {
	plib.RegisterItems("Swim", waterStuff)
}
