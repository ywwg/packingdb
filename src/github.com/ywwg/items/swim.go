package items

import (
	plib "github.com/ywwg/packinglib"
)

var waterStuff = []plib.Item{
	plib.NewBasicItem("froggies", []string{"Swimming"}, nil),
	plib.NewConsumableItem("swim suit", 0.25, plib.NoUnits, []string{"Swimming"}, nil),
	plib.NewBasicItem("swim towel", []string{"Swimming"}, nil),
	plib.NewBasicItem("plastic bag for wet things", []string{"Swimming", "Dirt"}, nil),
}

func init() {
	plib.RegisterItems("Swim", waterStuff)
}
