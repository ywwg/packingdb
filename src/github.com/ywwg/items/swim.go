package items

import (
	plib "github.com/ywwg/packinglib"
)

var waterStuff = []plib.Item{
	plib.NewBasicItem("froggies", []string{"Swimming"}, nil),
	plib.NewConsumableItem("swim suit", 0.5, plib.NoUnits, []string{"Swimming"}, nil),
	plib.NewBasicItem("swim towel", []string{"Swimming"}, nil),
}

func init() {
	plib.RegisterItems("Swim Stuff", waterStuff)
}
