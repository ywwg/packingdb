package items

import (
	plib "github.com/ywwg/packingdb/pkg/packinglib"
)

var waterStuff = []*plib.Item{
	plib.NewItem("swim suit", []string{"Swimming", "Lodging"}, nil).Consumable(0.25).Max(2.0),
	plib.NewItem("swim towel", []string{"Swimming"}, nil),
	plib.NewItem("drybag", []string{"Swimming", "Dirt"}, nil),
}

var tubing = []*plib.Item{
	plib.NewItem("innertube", []string{"Tubing"}, nil),
	plib.NewItem("air pump", []string{"Tubing"}, nil),
}
