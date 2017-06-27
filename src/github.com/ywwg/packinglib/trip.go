package packinglib

import (
//	"fmt"
)

// Trip describes a trip, which includes a length and a context
type Trip struct {
	Days int
	C    *Context
}

// MakeList returns a slice of PackedItems for the given trip
func (t *Trip) MakeList() []Item {
	AllItems := Clothing
	//	AllItems = append(AllItems, CampStuff...)
	//	AllItems = append(AllItems, WaterStuff...)
	//	AllItems = append(AllItems, Nevermore...)

	var packed []Item
	for _, i := range AllItems {
		p := i.Pack(t)
		if p.Count() > 0 {
			packed = append(packed, p)
		}
	}

	return packed
}
